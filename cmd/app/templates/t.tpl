import sdkit
from sdkit.models import load_model
from sdkit.generate import generate_images
from sdkit.utils import log

context = sdkit.Context()
context.device: str = 'cpu'
context.num_inference_steps: int = {{.Steps}}
context.sampler: str = '{{.Sampler}}'

# set the path to the model file on the disk (.ckpt or .safetensors file)
context.model_paths['stable-diffusion'] = '/home/bjones/easy-diffusion/models/stable-diffusion/{{.Model}}'
load_model(context, 'stable-diffusion')

# generate the image

images = generate_images(context, prompt='{{ .Prompt }}', negative_prompt='{{.NPrompt}}', seed={{.RandomNumber}}, width={{.Width}} , height={{.Height}},num_inference_steps={{.Steps}})

# save the image
images[0].save("{{.SmallImage}}") # images is a list of PIL.Image

#log.info("Generated images!")