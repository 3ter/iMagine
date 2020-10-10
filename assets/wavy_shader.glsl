#version 330 core
// The first line in glsl source code must always start with a version directive as seen above.
// GLSL stands for GL Shader Language

// vTexCoords are the texture coordinates, provided by Pixel
in vec2  vTexCoords;

// fragColor is what the final result is and will be rendered to your screen.
out vec4 fragColor;

// uTexBounds is the texture's boundaries, provided by Pixel.
uniform vec4 uTexBounds;

// uTexture is the actualy texture we are sampling from, also provided by Pixel.
uniform sampler2D uTexture;

// custom uniforms
uniform float uSpeed;
uniform float uTime;

void main() {
    // following this PR:
    // https://github.com/faiface/pixel-examples/pull/19/commits/6be8878d2a40e4afac719cbf15c596bc5b3e929c
	// vec2 t = gl_FragCoord.xy / uTexBounds.zw;
    vec2 t = vTexCoords / uTexBounds.zw;
	vec3 influence = texture(uTexture, t).rgb;

    if (influence.r + influence.g + influence.b > 0.3) {
		t.y += cos(t.x * 40.0 + (uTime * uSpeed))*0.005;
		t.x += cos(t.y * 40.0 + (uTime * uSpeed))*0.01;
	}

    vec3 col = texture(uTexture, t).rgb;
	fragColor = vec4(col * vec3(0.6, 0.6, 1.2),1.0);
}