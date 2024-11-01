from fastapi import APIRouter, status, Depends

from app.api.user.controllers import create, login
from app.models.user import UserCreateOutput, UserLoginOutput




router = APIRouter(prefix="/api/v1",tags=["user"])

router.add_api_route(path="/create",
                     endpoint=create,
                     status_code=status.HTTP_200_OK,
                     response_model=UserCreateOutput
                     methods=["POST"]
                    )

router.add_api_route(path="/login",
                     endpoint=login,
                     status_code=status.HTTP_200_OK,
                     response_model=UserLoginOutput,
                     methods=["POST"],
                    )



