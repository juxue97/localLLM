from fastapi import APIRouter, Depends, status

from app.api.chatbot.controllers import localLLM
from app.middleware.auth import authentication


router = APIRouter(prefix="/api/v1",tags=["localLLM"])




router.add_api_route(path="/localLLM",
                     endpoint=localLLM, #"#TODO endpoint function",
                     status_code=status.HTTP_200_OK,
                     methods=['POST'],
                     dependencies=[
                         Depends(authentication)
                    ]
                )