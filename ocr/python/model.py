from torch import nn
import torch.nn.functional as F


class OCRModel(nn.Module):
    def __init__(self):
        super(OCRModel, self).__init__()
        # input : BS x 1 x 64 x 64
        self.res_block1 = self._make_res_block(1, 16, 3, 1, 1)
        # input : BS x 16 x 64 x 64
        self.max_pool1 = nn.MaxPool2d(2, 2)
        # input : BS x 16 x 32 x 32
        self.res_block2 = self._make_res_block(16, 32, 3, 1, 1)
        # input : BS x 32 x 32 x 32
        self.max_pool2 = nn.MaxPool2d(2, 2)
        # input : BS x 32 x 16 x 16
        self.res_block3 = self._make_res_block(32, 64, 3, 1, 1)
        # input : BS x 64 x 16 x 16
        self.max_pool3 = nn.MaxPool2d(2, 2)
        # input : BS x 64 x 8 x 8
        self.fc1 = nn.Linear(64 * 8 * 8, 100)
        self.fc2 = nn.Linear(100, 10)

        self.dropout = nn.Dropout(0.5)

    def forward(self, x):
        x = self.res_block1(x)
        x = self.max_pool1(x)
        x = self.res_block2(x)
        x = self.max_pool2(x)
        x = self.res_block3(x)
        x = self.max_pool3(x)
        x = nn.Flatten()(x)
        x = self.fc1(x)
        x = F.relu(x)
        x = self.dropout(x)
        x = self.fc2(x)

        return x

    def _make_res_block(self, in_channels, out_channels, kernel_size, stride, padding):
        return nn.Sequential(
            nn.Conv2d(in_channels, out_channels, kernel_size, stride, padding),
            nn.BatchNorm2d(out_channels),
            nn.ReLU(),
            nn.Conv2d(out_channels, out_channels, kernel_size, stride, padding),
            nn.BatchNorm2d(out_channels),
            nn.ReLU(),
        )
