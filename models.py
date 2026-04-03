from typing import Optional, List
from pydantic import BaseModel, field_validator
from sqlmodel import Field, SQLModel, Relationship, create_engine

class CPTagLink(SQLModel, table=True):
    # 注意：这里的 foreign_key 必须指向 cp 表的 id
    cp_id: Optional[int] = Field(
        default=None, foreign_key="cp.id", primary_key=True
    )
    # 注意：这里的 foreign_key 必须指向 tag 表的 id
    tag_id: Optional[int] = Field(
        default=None, foreign_key="tag.id", primary_key=True
    )

class CPBase(SQLModel):
    name: str = Field(index=True, min_length=1, unique=True)  # 强制校验：至少1个字符
    category: str = Field(min_length=2, max_length=20) # 强制校验：2-10个字符
    link: Optional[str] = Field(default=None, nullable=True)

class TagRead(SQLModel):
    id: int
    name: str

class CPRead(CPBase):
    id: int
    tags: List[TagRead] = []

class CPResponse(BaseModel):
    status: str
    data: List[CPRead] | CPRead

# 1. 专门用于接收前端数据的 Schema (不对应数据库表)
class CPCreate(CPBase):
    tag_names: List[str] = Field(min_length=1) # 这个字段只存在于内存，用于接收 POST 数据

    @field_validator("tag_names")
    @classmethod
    def check_tags_not_empty(cls, v: List[str]) -> List[str]:
        if not v or len([name for name in v if name.strip()]) == 0:
            raise ValueError("At least one non-empty tag is required")
        return v

# 2. 真正的数据库表模型
class CP(CPBase, table=True):
    id: Optional[int] = Field(default=None, primary_key=True)

    # 关系 (这才是 SQLAlchemy 认识的 List)
    tags: List["Tag"] = Relationship(back_populates="cps", link_model=CPTagLink)



class Tag(SQLModel, table=True):
    id: Optional[int] = Field(default=None, primary_key=True)
    name: str = Field(index=True, unique=True) # 标签名通常是唯一的
    
    # 通过中间表关联到文章
    cps: List["CP"] = Relationship(
        back_populates="tags", link_model=CPTagLink
    )

# 1. 定义数据库连接地址
# 使用 SQLite 并在当前目录下生成名为 database.db 的文件
sqlite_file_name = "database.db"
sqlite_url = f"sqlite:///{sqlite_file_name}"

# 2. 创建引擎 (Engine)
# echo=True 会在终端打印出生成的 SQL 语句，非常适合调试查看表结构
engine = create_engine(sqlite_url, echo=True)

# 3. 执行创建动作
def create_db_and_tables():
    # 这行代码会扫描所有继承自 SQLModel 且 table=True 的类
    # 并根据它们的定义在数据库中创建相应的表
    SQLModel.metadata.create_all(engine)

if __name__ == "__main__":
    create_db_and_tables()
    print("数据库表创建成功！")