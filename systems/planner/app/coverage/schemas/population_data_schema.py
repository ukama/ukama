from sqlalchemy import Column, Float, Integer, String, Boolean
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker


# SQLAlchemy Configuration
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
