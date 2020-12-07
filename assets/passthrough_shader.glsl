#version 330 core
//This shader is meant to pass all values through as provided by Pixel
//copied from the grayscale shader pixel example
//TODO currently this seems to set everything to grayscale, which isn't noticeable in the demo, but in other scenes

in vec2  vTexCoords;

out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

void main() {
	// Get our current screen coordinate
	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;

	vec4 color = vec4( texture(uTexture, t).rgb, 1.0);
	
	fragColor = color;
}

