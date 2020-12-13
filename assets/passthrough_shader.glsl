#version 330 core
//This shader is meant to pass all values through as provided by Pixel
//TODO currently this seems to set everything to grayscale, which isn't noticeable in the demo, but in other scenes

//This is provided by the vertex shader(?), the variable needs to be this name in order 
//for the shader to work
in vec2  vTexCoords;


//This is the output required of the fragment shader, the variable can be named anything we want, it just needs to be a vec4
out vec4 fragColor;


uniform vec4 uTexBounds;
uniform sampler2D uTexture;

void main() {
	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;

	vec4 color = vec4( texture(uTexture, t).rgb, 0.9);
	
	fragColor = color;
}

