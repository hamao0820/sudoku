import torch


# 画像の四隅にランダムで黒い線を引く
class RandomLine(torch.nn.Module):
    def __init__(self, p=0.5, thickness=1):
        super(RandomLine, self).__init__()
        self.p = p
        self.thickness = thickness

    def forward(self, x):
        # 線の太さ
        # 右端の線
        if torch.rand(1) < self.p:
            x[:, :, -self.thickness :] = 0
        # 下端の線
        if torch.rand(1) < self.p:
            x[:, -self.thickness :, :] = 0
        # 左端の線
        if torch.rand(1) < self.p:
            x[:, :, : self.thickness] = 0
        # 上端の線
        if torch.rand(1) < self.p:
            x[:, : self.thickness, :] = 0

        return x


if __name__ == "__main__":
    import torchvision.transforms as transforms
    from PIL import Image

    img = Image.open("ocr/data/train/train/1/1.png")
    img.show()
    trans = transforms.Compose([transforms.ToTensor(), RandomLine(p=0.5), transforms.ToPILImage()])
    trans(img).show()
