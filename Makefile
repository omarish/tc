VERSION ?= dev

# Demo output is 16:9 for social embeds; vhs renders tight and we pad to fit.
DEMO_VER ?= v0.2.0

.PHONY: build test clean demo

build:
	go build -ldflags "-X main.version=$(VERSION)" -o tc .

test:
	go test ./...

clean:
	rm -f tc

# Render the release demo (GIF + MP4) from demo/tc.tape.
# Requires vhs: brew install vhs
demo: build
	@PATH="$(CURDIR):$$PATH" vhs demo/tc.tape
	@# Viewport already fits the content, so no padding pass is needed.
	@ffmpeg -v error -i demo/.scratch.mp4 -vf fps=30 \
	  -c:v libx264 -profile:v high -pix_fmt yuv420p -crf 20 -movflags +faststart \
	  demo/tc-$(DEMO_VER).mp4 -y
	@rm -f demo/.scratch.mp4
	@ffprobe -v error -select_streams v:0 -show_entries stream=width,height,r_frame_rate \
	  -show_entries format=duration -of default=nw=1 demo/tc-$(DEMO_VER).mp4