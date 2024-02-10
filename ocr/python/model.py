from torch import nn
from torchvision import models
import torch.nn.functional as F


class OCRModel(nn.Module):
    def __init__(self):
        super(OCRModel, self).__init__()
        self.resnet18 = models.resnet18(pretrained=True)
        for param in self.resnet18.parameters():
            param.requires_grad = False

        self.resnet18.fc = nn.Linear(512, 10)

    def forward(self, x):
        x = self.resnet18(x)
        return F.log_softmax(x, dim=1)
