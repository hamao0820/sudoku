import random

import torch
from torch.utils.data import DataLoader
from torchvision import transforms
from torch import nn, optim
from tqdm import tqdm

from dataset import OCRDataset, Subset
from model import OCRModel


def train():
    # device
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")

    # dataset
    common_transforms = transforms.Compose([transforms.Resize((64, 64))])
    test_transforms = transforms.Compose(
        [
            transforms.ToTensor(),
            transforms.GaussianBlur(3),
            transforms.Normalize((0.5,), (0.5,)),
        ]
    )
    augment_transform = transforms.Compose(
        [
            transforms.RandomOrder(
                [
                    transforms.RandomResizedCrop(64, scale=(0.09, 1.0)),
                    transforms.RandomChoice(
                        [
                            transforms.RandomAffine(degrees=10, translate=(0.1, 0.1), scale=(0.9, 1.1), shear=15),
                            transforms.RandomRotation(10),
                        ]
                    ),
                    transforms.RandomPerspective(distortion_scale=0.2, p=0.5),
                    transforms.RandomPosterize(bits=4),
                ]
            ),
            transforms.GaussianBlur(3),
            transforms.ToTensor(),
            transforms.Normalize((0.5,), (0.5,)),
            transforms.RandomErasing(p=0.5, scale=(0.01, 0.1), value=0),
        ]
    )

    ds = OCRDataset(common_transforms)
    indeces = list(range(len(ds)))
    train_size = int(0.8 * len(ds))
    random.shuffle(indeces)
    train_ds = Subset(ds, indeces[:train_size], augment_transform)
    val_ds = Subset(ds, indeces[train_size:], test_transforms)
    train_dl = DataLoader(train_ds, batch_size=128, shuffle=True)
    val_dl = DataLoader(val_ds, batch_size=128, shuffle=True)

    # model
    model = OCRModel()
    model.to(device)

    # loss and optimizer
    criterion = nn.CrossEntropyLoss()
    optimizer = optim.Adam(model.parameters(), lr=0.0001)

    # train
    corrects = []
    for epoch in range(200):
        print(f"Epoch {epoch + 1}")
        for i, data in tqdm(enumerate(train_dl, 0), total=len(train_dl)):
            inputs, labels = data
            inputs, labels = inputs.to(device), labels.to(device)

            optimizer.zero_grad()

            outputs = model(inputs)
            loss = criterion(outputs, labels)
            loss.backward()
            optimizer.step()

        # validation
        model.eval()
        correct = 0
        total = 0
        with torch.no_grad():
            for data in val_dl:
                images, labels = data
                images, labels = images.to(device), labels.to(device)
                outputs = model(images)
                _, predicted = torch.max(outputs.data, 1)
                total += labels.size(0)
                # output and labels are one-hot encoded
                correct += (predicted == torch.argmax(labels, 1)).sum().item()
        tqdm.write(f"val acc: {100 * correct / total:.2f}%")
        corrects.append(correct / total)
        model.train()

        if epoch % 10 == 0:
            torch.save(model.state_dict(), "ocr/python/model/model2.pth")

    torch.save(model.state_dict(), "ocr/python/model/model2.pth")

    with open("ocr/python/corrects2.txt", "w") as f:
        f.write("\n".join([str(c) for c in corrects]))

    print("Finished Training")


if __name__ == "__main__":
    train()
