import torch
from torch.utils.data import DataLoader
from torchvision import transforms
from torch import nn, optim
from tqdm import tqdm
from timm.scheduler import CosineLRScheduler

from dataset import OCRDataset
from model import OCRModel
from transform import RandomLine


def train():
    # device
    device = torch.device("mps" if torch.backends.mps.is_available() else "cpu")
    print(f"device: {device}")

    # dataset
    test_transforms = transforms.Compose(
        [
            transforms.Resize((224, 224)),
            transforms.ToTensor(),
        ]
    )
    augment_transform = transforms.Compose(
        [
            RandomLine(p=0.2),
            transforms.Resize((224, 224)),
            transforms.RandomChoice(
                [
                    transforms.RandomAffine(degrees=5, translate=(0.1, 0.1), scale=(0.9, 1.1), shear=5),
                    transforms.RandomPerspective(distortion_scale=0.2, p=0.5),
                ]
            ),
            transforms.RandomResizedCrop(64, scale=(0.8, 1.0), ratio=(0.9, 1.1)),
            transforms.ToTensor(),
        ],
    )

    train_ds = OCRDataset("train", transform=augment_transform)
    val_ds = OCRDataset("valid", transform=test_transforms)
    train_dl = DataLoader(train_ds, batch_size=64, shuffle=True)
    val_dl = DataLoader(val_ds, batch_size=64, shuffle=True)

    # model
    model = OCRModel()
    model.to(device)

    # loss and optimizer
    criterion = nn.CrossEntropyLoss()
    optimizer = optim.Adam(model.parameters(), lr=0.0001)
    scheduler = CosineLRScheduler(
        optimizer, t_initial=100, lr_min=1e-6, warmup_t=3, warmup_lr_init=1e-6, warmup_prefix=True
    )

    # train
    corrects = []
    for epoch in range(100):
        print(f"Epoch {epoch + 1}")
        for i, data in tqdm(enumerate(train_dl, 0), total=len(train_dl)):
            inputs, labels = data
            inputs, labels = inputs.to(device), labels.to(device)

            optimizer.zero_grad()

            outputs = model(inputs)
            loss = criterion(outputs, labels)
            loss.backward()

            optimizer.step()

        scheduler.step(epoch)

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

        if epoch % 5 == 0:
            torch.save(model.state_dict(), "ocr/python/model/model.pth")
            print(f"epoch: {epoch + 1}, model saved")

    torch.save(model.state_dict(), "ocr/python/model/model.pth")
    with open("ocr/python/corrects.txt", "w") as f:
        f.write("\n".join([str(c) for c in corrects]))

    print("Finished Training")


if __name__ == "__main__":
    train()
