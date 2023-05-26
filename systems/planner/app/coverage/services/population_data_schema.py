from sqlalchemy import create_engine, Column, Float, Integer, String, Boolean
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker


# SQLAlchemy Configuration
SQLALCHEMY_DATABASE_URL = "mysql+mysqlconnector://root:MyNewPass@localhost/planner_tool"

engine = create_engine(SQLALCHEMY_DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base = declarative_base()


class PopulationData(Base):
    __tablename__ = 'population_data_simple'

    id = Column(Integer, primary_key=True, index=True)
    longitude = Column(Float)
    latitude = Column(Float)
    value = Column(Float)
    

class FilesStatus(Base):
    __tablename__ = 'files_status'

    id = Column(Integer, primary_key=True, index=True)
    file_path = Column(String(255))
    parsed = Column(Boolean)


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()