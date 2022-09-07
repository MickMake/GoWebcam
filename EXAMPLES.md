# Examples

## Running examples

# Simple cron pulling an image every 5 minutes.
	% GoWebcam cron run  . ./5 . . . .  webcam get Basin https://charlottepass.com.au/charlottepass/webcam/lucylodge/current.jpg

# Once-off run of all webcams defined in config.json file.
	% GoWebcam webcam run

# Pull webcam images as defined in config.json file, via cron.
	% GoWebcam webcam cron

# Same as above, but run as a daemon.
	% GoWebcam daemon exec  webcam cron

# List currently scheduled jobs.
	% GoWebcam daemon list
