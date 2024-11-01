from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.api.chatbot.routes import router

app = FastAPI()



origins = ["*"]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/")
async def root():
    return {"HEALTH_CHECK":"OK", "CONNECTION":"OK"}


app.include_router(router)