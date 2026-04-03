from typing import Any, Dict
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
from sqlmodel import col, create_engine, Session, select
from sqlalchemy.orm import selectinload
from models import CP, CPCreate, CPRead, Tag, CPResponse

# 1. 定义数据库连接地址
# 使用 SQLite 并在当前目录下生成名为 database.db 的文件
sqlite_file_name = "database.db"
sqlite_url = f"sqlite:///{sqlite_file_name}"

# 2. 创建引擎 (Engine)
# echo=True 会在终端打印出生成的 SQL 语句，非常适合调试查看表结构
engine = create_engine(sqlite_url, echo=True)

app = FastAPI()

class AppError(Exception):
    def __init__(
        self, 
        message: str, 
        status_code: int = 400, 
        **kwargs: Any
    ) -> None:
        self.message = message
        self.status_code = status_code
        self.extra = kwargs  # 接收任何额外的关键字参数
        super().__init__(message)

@app.exception_handler(AppError)
async def universal_exception_handler(request: Request, exc: AppError):
    # 直接取属性，不要去算 len(exc.args)
    return JSONResponse(
        status_code=exc.status_code,
        content={
            "status": "error",
            "msg": exc.message,
            **exc.extra  # 所有的自定义参数会自动解包到这里
        }
    )


@app.get("/")
async def root():
    return {"message": "Hello World"}

@app.get("/cp", response_model=CPResponse)
async def get_all_cp() -> Dict[str, Any]:
    with Session(engine) as session:
        statement = select(CP).options(selectinload(getattr(CP, "tags")))

        result = session.exec(statement)

        cps = result.all()
        return {"status": "ok", "data": cps}

@app.post("/cp", response_model=CPResponse)
async def create_cp(cp_in: CPCreate) -> Dict[str, Any]:
    with Session(engine) as session:
        # 1. 检查 CP 是否重名
        existing_cp = session.exec(select(CP).where(col(CP.name) == cp_in.name)).first()
        if existing_cp:
            raise AppError("The CP name already exists.", 400)

        # 2. 将输入转换为数据库模型 (注意：此时 tags 列表还是空的)
        db_cp = CP.model_validate(cp_in)

        # 3. 处理标签逻辑
        if cp_in.tag_names:
            # 去重：防止前端传了重复的标签名
            unique_tag_names = list(set(cp_in.tag_names))
            
            # 一次性查出数据库中已有的标签对象
            statement = select(Tag).where(col(Tag.name).in_(unique_tag_names))
            existing_tags = session.exec(statement).all()
            
            # 找出哪些标签是数据库里没有的，并创建它们
            existing_names = {t.name for t in existing_tags}
            new_tags = [Tag(name=name) for name in unique_tag_names if name not in existing_names]
            
            # 合并：已有的 + 新建的
            db_cp.tags = [*existing_tags, *new_tags]

        # 4. 提交到数据库
        # SQLModel 会自动处理：1. 插入 CP 2. 插入新 Tag 3. 在 CPTagLink 插入关联
        session.add(db_cp)
        session.commit()
        session.refresh(db_cp)

        _ = db_cp.tags

        return {"status": "ok", "data": db_cp}

@app.delete("/cp/{cp_id}", response_model=CPResponse)
async def delete_cp(cp_id: int) -> Dict[str, Any]:
    with Session(engine) as session:
        statement = select(CP).options(selectinload(getattr(CP, "tags"))).where(CP.id == cp_id)
        db_cp = session.exec(statement).first()
        if not db_cp:
            raise AppError("CP not found", 404)
        
        response_data = CPRead.model_validate(db_cp)
        
        # 3. 执行真实的删除和提交操作
        session.delete(db_cp)
        session.commit()
        
    # 4. 返回的不再是数据库对象，而是刚刚提取出来的纯数据对象
    return {"status": "ok", "data": response_data}
@app.put("/cp/{cp_id}", response_model=CPResponse)
async def update_cp(cp_id: int, cp_in: CPCreate) -> Dict[str, Any]:
    with Session(engine) as session:
        # 1. 查找旧数据
        db_cp = session.get(CP, cp_id)
        if not db_cp:
            raise AppError("CP not found", 404)
        
        # 2. 更新基础字段
        db_cp.name = cp_in.name
        db_cp.category = cp_in.category
        db_cp.link = cp_in.link

        # 3. 更新标签 (核心逻辑)
        unique_tag_names = list(set(cp_in.tag_names))
        existing_tags = session.exec(select(Tag).where(col(Tag.name).in_(unique_tag_names))).all()
        existing_names = {t.name for t in existing_tags}
        new_tags = [Tag(name=name) for name in unique_tag_names if name not in existing_names]
        
        # 直接赋值！SQLAlchemy 会自动计算 Diff：
        # - 消失的标签：自动从中间表删掉关联
        # - 新增的标签：自动在中间表加关联
        db_cp.tags = [*existing_tags, *new_tags]

        session.add(db_cp)
        session.commit()
        session.refresh(db_cp)

        _ = db_cp.tags

        return {"status": "ok", "data": db_cp}