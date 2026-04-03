from typing import Any, Dict
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
from sqlmodel import create_engine, Session, select
from models import CP, CPBase

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

@app.get("/cp")
async def get_all_cp() -> Dict[str, Any]:
    with Session(engine) as session:
        statement = select(CP)

        result = session.exec(statement)

        cps = result.all()
        return {"status": "ok", "data": cps}

@app.post("/cp")
async def create_cp(cp: CPBase) -> Dict[str, Any]:
    cp = CP.model_validate(cp)
    with Session(engine) as session:
        existing_cp = session.exec(select(CP).where(CP.name == cp.name)).first()
        if existing_cp:
            raise AppError("The CP name already exists.", 400)
        session.add(cp)
        session.commit()
        session.refresh(cp)
    return {"status": "ok", "data": cp}

@app.delete("/cp/{cp_id}")
async def delete_cp(cp_id: int) -> Dict[str, Any]:
    with Session(engine) as session:
        statement = select(CP).where(CP.id == cp_id)
        cp = session.exec(statement).first()
        if not cp:
            raise AppError("CP not found", 404)
        session.delete(cp)
        session.commit()
    return {"status": "ok", "data": cp}

@app.put("/cp/{cp_id}")
async def update_cp(cp_id: int, cp: CPBase) -> Dict[str, Any]:
    cp = CP.model_validate(cp)
    with Session(engine) as session:
        statement = select(CP).where(CP.id == cp_id)
        existing_cp = session.exec(statement).first()
        if not existing_cp:
            raise AppError("CP not found", 404)
        
        existing_cp.name = cp.name
        existing_cp.category = cp.category
        existing_cp.link = cp.link

        session.add(existing_cp)
        session.commit()
        session.refresh(existing_cp)
    return {"status": "ok", "data": existing_cp}