import sdkit
from sdkit.models import load_model
from sdkit.generate import generate_images
from sdkit.utils import log

context = sdkit.Context()
context.device: str = 'cpu'
context.num_inference_steps: int = 16

# set the path to the model file on the disk (.ckpt or .safetensors file)
context.model_paths['stable-diffusion'] = '/home/bjones/easy-diffusion/models/stable-diffusion/{{.Model}}'
load_model(context, 'stable-diffusion')

# generate the image

images = generate_images(context, prompt='{{ .Prompt }}', negative_prompt='{{.NPrompt}}', seed={{.RandomNumber}}, width=512 , height=768,num_inference_steps=16)

# save the image
images[0].save("{{.SmallImage}}") # images is a list of PIL.Image

#log.info("Generated images!")