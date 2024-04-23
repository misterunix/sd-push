from PIL import Image

import sdkit
from sdkit.filter import apply_filters
from sdkit.models import load_model

context = sdkit.Context()
context.device: str = 'cpu'
image = Image.open("{{ .SmallImage }}")

# set the path to the model file on the disk
context.model_paths["realesrgan"] = "models/realesrgan/RealESRGAN_x4plus.pth"
load_model(context, "realesrgan")

# apply the filter
scale = 4  # or 2
image_upscaled = apply_filters(context, "realesrgan", image, scale=scale)

# save the filtered image
image_upscaled[0].save("{{ .LargeImage }}")

