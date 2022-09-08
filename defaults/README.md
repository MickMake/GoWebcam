# GoWebcam - A simple automated webcam fetcher written in GoLang


## What does it do?
This is a simple webcam fetcher. It was an itch I needed to scratch.
What it does is:
1. Regularly pull webcam images from any URL, (supports authentication).
2. Only creates a new file if the image has changed.
3. Allows for renaming of files based on time rounding, or Tesseract OCR.
4. Run scripts periodically. EG: To create MP4 videos from captured images.
5. Simple JSON config file.


## Running examples

### Simple cron pulling an image every 5 minutes.
	% GoWebcam cron run  . ./5 . . . .  webcam get Basin https://charlottepass.com.au/charlottepass/webcam/lucylodge/current.jpg

### Once-off run of all webcams defined in config.json file.
	% GoWebcam webcam run

### Pull webcam images as defined in config.json file, via cron.
	% GoWebcam webcam cron

### Same as above, but run as a daemon.
	% GoWebcam daemon exec  webcam cron

### List currently scheduled jobs.
	% GoWebcam daemon list


## Configuration - config.json

The config file is broken up into sections.


### "images"
Defines webcam images to fetch.

Requires:
- `prefix` - The name of the webcam. Used as a prefix for directory and file names.
- `url` - The webcam url to fetch.

Optional:
- `cron` - The schedule to fetch. Will default to the global `cron` definition.
- `timeout` - Defines the web fetch timeout.
- `dir` - Alternative based directory.

```
{
	"images": [
		{
			"prefix": "Basin",
			"url": "https://charlottepass.com.au/charlottepass/webcam/lucylodge/current.jpg",
			"cron": "30 */5 * * * *",
	        "timeout": 30000000,
	        "dir": "images"
		}
}
```

### Globals
Global options.

Optional:
- `cron` - Define global cron for webcam image fetches. Default: `00 */5 * * * *`.
- `timeout` - Defines the web fetch timeout. Default: `30s`.
- `dir` - Alternative based directory. Default: `images`.

```
{
	"timeout": 30000000,
	"dir": "images",
	"cron": "30 */5 * * * *"
}
```

### "rename"
Rename files based on rules.

Optional:
- `by_time` - Round down timestamp filename to nearest duration.
- `ocr` - Not yet supported.

```
{
	"rename": {
		"by_time": "5m",
		"ocr": ""
	}
}
```

### "report"
Regular scheduled reporting.

Required:
- `cron` - The schedule to produce reports.

Optional:
- `level` - Not yet supported. Default: just print running jobs.

```
{
	"report": {
		"level": "jobs",
		"cron": "00 00 */1 * * *"
	}
}
```

### "scripts"
Schedule scripts to run.

Required:
- `cron` - The schedule to run the script.
- `cmd` - The script to run.

Optional:
- `args` - Command arguments.

```
{
	"scripts": [
		{
			"cron": "00 01 21 * * *",
			"cmd": "./process_yearly.sh",
			"args": []
		}
	]
}
```


## Example config.json file

```
{
	"images": [
		{
			"prefix": "Basin",
			"url": "https://charlottepass.com.au/charlottepass/webcam/lucylodge/current.jpg"
		},
		{
			"prefix": "KangarooRidge",
			"url": "https://charlottepass.com.au/charlottepass/webcam/Chairlift/current.jpg"
		},
		{
			"prefix": "Pulpit",
			"url": "https://charlottepass.com.au/charlottepass/webcam/timber/current.jpg"
		},
		{
			"prefix": "Guthries",
			"url": "https://charlottepass.com.au/charlottepass/webcam/guthries/current.jpg",
			"cron": "00 */10 * * * *"
		}
	],
	"timeout": 30000000,
	"dir": "images",
	"cron": "30 */5 * * * *",
	"rename": {
		"by_time": "5m"
	},
	"report": {
		"level": "jobs",
		"cron": "00 00 */1 * * *"
	},
	"scripts": [
		{
			"cron": "00 01 21 * * *",
			"cmd": "./process_yearly_tl.sh",
			"args": []
		},
		{
			"cron": "00 56 */1 * * *",
			"cmd": "./process_daily_tl.sh",
			"args": []
		},
		{
			"cron": "30 00 */2 * * *",
			"cmd": "uname",
			"args": ["-a"]
		}
	]
}

```
