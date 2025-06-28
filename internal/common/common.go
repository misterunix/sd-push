package common

import "os/exec"

type Startup struct {
	RandomNumber   int
	Prompt         string
	NPrompt        string
	ModelsLocation string
	Model          string
	LoraLocation   string
	Steps          int
	Width          int
	Height         int
	ScaleUp        bool
	RemoveSmall    bool
	Sampler        string
	Count          int
}

type Stable struct {
	RandomNumber int
	SmallImage   string
	LargeImage   string
	Prompt       string
	NPrompt      string
	Model        string
	Steps        int
	Width        int
	Height       int
	Sampler      string
}

var Samplers []string     // List of samplers
var BaseModels []string   // List of models in the directory
var RunThisSampler string // The sampler to run, set by the user

// var timedir string
var Tss string
var Cmd *exec.Cmd
var Width, Height int
var UserDir string

/*
type EasyDiffusion struct {
  Prompt string // The text prompt to generate the image
  Negative_prompt string // The negative prompt to avoid certain features in the image
  Seed int // The random seed for image generation
  Used_random_seed bool // Indicates if a random seed was used
  Num_outputs int // The number of images to generate, with these parameters, it is set to 1
  Num_inference_steps int // The number of inference steps for the diffusion model
  Guidance_scale float64 // The guidance scale for the model, which influences the adherence to the prompt
  Width int // The width of the generated image
  Height int // The height of the generated image
  Vram_usage_level string // The VRAM usage level, set to "high" for better quality
  Sampler_name string // The name of the sampler to use, set to "k_euler_ancestral" for this example
  Use_stable_diffusion_model string // The model to use for image generation, set to "sd15" for this example
  Clip_skip bool // Whether to use clip skip, set to false in this example
  Use_vae_model string // The VAE model to use, set to "vae-ft-mse-840000-ema-pruned" for this example
  Stream_progress_updates bool // Whether to stream progress updates, set to false in this example
  "stream_image_progress": false,
  "show_only_filtered_image": true,
  "block_nsfw": false,
  "output_format": "png",
  "output_quality": 75,
  "output_lossless": false,
  "metadata_output_format": "json",
  "original_prompt": "Masterpiece, extreme contrast, brushstrokes, very strong highlights, 'standing in a torrent of ink' pose, abstract and mysterious background (semi-transparent white ink painting: 1.5), (rough and thick brushstrokes: 1.5), semi-realistic, (art filter: 1.5), (breathtakingly beautiful Japanese woman: 1.5), modern Japanese, Using medium green and pastel orange as accent colors, the figure is monochrome",
  "active_tags": [],
  "inactive_tags": [],
  "use_upscale": "RealESRGAN_x4plus",
  "upscale_amount": "4",
  "enable_vae_tiling": true
*/
