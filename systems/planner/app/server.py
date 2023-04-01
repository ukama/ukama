from fastapi import FastAPI, Depends

from api import router
from api.v1.coverage import coverage_router
from api.home.home import home_router
from core.config import config
from core.fastapi.dependencies import Logging


def init_routers(app_: FastAPI) -> None:
    app_.include_router(home_router)
    app_.include_router(router)

def create_app() -> FastAPI:
    app_ = FastAPI(
        title="Planner System",
        description="Planner System APIs",
        version="1.0.0",
        docs_url=None if config.ENV == "production" else "/docs",
        redoc_url=None if config.ENV == "production" else "/redoc",
        dependencies=[Depends(Logging)],
    )
    init_routers(app_=app_)
    return app_


app = create_app()
