from glob import glob
import os

import torch
from torch.utils.data import Dataset
from torchvision import transforms
from PIL import Image


class OCRDataset(Dataset):
    def __init__(self, transform=None):
        self.data = []
        dirs = glob("ocr/data/cells/**")
        for d in dirs:
            label = os.path.basename(d)
            files = glob(d + "/*.png")
            for f in files:
                self.data.append({"label": label, "path": f})
        self.transform = transform

    def __len__(self) -> int:
        return len(self.data)

    def __getitem__(self, idx) -> (torch.Tensor, int):
        item = self.data[idx]
        img = Image.open(item["path"])
        if self.transform:
            img = self.transform(img)
        # label to one-hot encoding
        label = torch.zeros(10)
        label[int(item["label"])] = 1
        return img, label


class Subset(OCRDataset):
    def __init__(self, dataset, indices, transform=None):
        self.data = dataset
        self.transform = transform
        self.indices = indices

    def __len__(self) -> int:
        return len(self.indices)

    def __getitem__(self, idx) -> (torch.Tensor, int):
        img, label = self.data[self.indices[idx]]
        if self.transform:
            img = self.transform(img)
        return img, label


# test
if __name__ == "__main__":
    ds = OCRDataset(transforms.ToTensor())
    print(len(ds))
    img, label = ds[0]
    print(label)
    print(img.shape)
