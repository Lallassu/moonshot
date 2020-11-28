#version 330 core
layout (location = 0) in vec3 pos;
layout (location = 1) in vec3 color;
layout (location = 2) in vec2 texCoord;

uniform mat4 model, view, projection;
out vec3 fColor;
out vec2 fCoord;

void main()
{
	fColor = color;
	fCoord = vec2(texCoord.x, texCoord.y);
    gl_Position =  projection * view * model * vec4(pos, 1.0);
}
