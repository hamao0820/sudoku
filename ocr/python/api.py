from fastapi import FastAPI
from pydantic import BaseModel
import base64
from PIL import Image
from torchvision import transforms
from io import BytesIO
import torch


# load the model
model = torch.load("ocr/python/model/model.pt")
device = torch.device("mps" if torch.backends.mps.is_available() else "cpu")
model.to(device)
model.eval()

app = FastAPI()

trans = transforms.Compose(
    [
        transforms.Resize((64, 64)),
        transforms.ToTensor(),
    ]
)


class OCRReq(BaseModel):
    b64: str


@app.post("/ocr")
async def ocr(req: OCRReq):
    # convert the base64 to a PIL image
    img = base64_to_pil(req.b64)
    # convert the image to a tensor
    img = trans(img)
    # add a batch dimension
    img = img.unsqueeze(0)
    # send the tensor to the device
    img = img.to(device)
    # run the model
    with torch.no_grad():
        out = model(img)
    # get the predicted class
    pred = out.argmax(dim=1).item()
    # return the prediction
    return pred
    # return "hello"


def base64_to_pil(img_str):
    if "base64," in img_str:
        # DARA URI の場合、data:[<mediatype>][;base64], を除く
        img_str = img_str.split(",")[1]
    img_raw = base64.b64decode(img_str)
    img = Image.open(BytesIO(img_raw))

    return img


# run the app
if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        app,
        port=8888,
    )
