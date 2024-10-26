import os
from fastapi import APIRouter, FastAPI, Response, status
import uvicorn
from mcstatus import JavaServer
from dotenv import load_dotenv

# * for local development
load_dotenv()

server = JavaServer.lookup(os.getenv("MINECRAFT_SERVER_IP"))


app = FastAPI(openapi_url="/docs/openapi.json", docs_url="/docs")
# * ------------------------Root endpoint------------------------


# * required for fargate task group healthy check
@app.get("/", tags=["Root"])
async def read_root():
    return "Minecraft Bridge API"


# * ------------------------API section--------------------------
router = APIRouter(prefix="/api/v1/minecraft")


@router.get("/", tags=["Root"])
async def api_root():
    return "MCStatus API"


@router.get("/ping", tags=["mcstatus"])
async def ping(response: Response):
    try:
        latency = server.ping()
        response.status_code = status.HTTP_200_OK
        return {"latency": latency}
    except Exception as e:
        response.status_code = status.HTTP_500_INTERNAL_SERVER_ERROR
        return {"error": str(e)}


@router.get("/status", tags=["mcstatus"])
async def get_status(response: Response):
    try:
        status = server.status()
        response.status_code = status.HTTP_200_OK
        return {
            "players": {"online": status.players.online, "max": status.players.max},
            "version": status.version.name,
            "motd": status.motd.to_minecraft(),
            "icon": status.icon,
        }
    except Exception as e:
        response.status_code = status.HTTP_500_INTERNAL_SERVER_ERROR
        return {"error": str(e)}


app.include_router(router)

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=9000, reload=True)
