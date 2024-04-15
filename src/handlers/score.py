# Based on https://github.com/christophschuhmann/improved-aesthetic-predictor/blob/fe88a163f4661b4ddabba0751ff645e2e620746e/simple_inference.py

import io
from importlib import resources
from typing_extensions import assert_type

import torch
import torch.nn as nn
from transformers import CLIPModel, CLIPProcessor

from src.database.query import Querier


class MLP(nn.Module):
    def __init__(self):
        super().__init__()
        self.layers = nn.Sequential(
            nn.Linear(768, 1024),
            nn.Dropout(0.2),
            nn.Linear(1024, 128),
            nn.Dropout(0.2),
            nn.Linear(128, 64),
            nn.Dropout(0.1),
            nn.Linear(64, 16),
            nn.Linear(16, 1),
        )

    @torch.no_grad()
    def forward(self, embed):
        return self.layers(embed)


class AestheticScorer(torch.nn.Module):
    def __init__(self, dtype, querier: Querier):
        super().__init__()
        self.clip = CLIPModel.from_pretrained("openai/clip-vit-large-patch14")
        self.processor = CLIPProcessor.from_pretrained("openai/clip-vit-large-patch14")
        self.mlp = MLP()
        scorer = querier.get_scorer(name="sac+logos+ava1-l14-linearMSE.pth")
        state_dict = torch.load(io.BytesIO(scorer.state_dict.tobytes()))  # pyright: ignore
        self.mlp.load_state_dict(state_dict)
        self.dtype = dtype
        self.eval()

    @torch.no_grad()
    def __call__(self, images):
        device = next(self.parameters()).device
        inputs = self.processor(images=images, return_tensors="pt")  # pyright: ignore
        inputs = {k: v.to(self.dtype).to(device) for k, v in inputs.items()}
        embed = self.clip.get_image_features(**inputs)  # pyright: ignore
        # normalize embedding
        embed = embed / torch.linalg.vector_norm(embed, dim=-1, keepdim=True)
        return self.mlp(embed).squeeze(1)
