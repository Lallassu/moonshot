#version 330 core

uniform mat4 model, view, projection;

layout (location = 0) in vec3 pos;
layout (location = 1) in vec4 color;

out vec4 fragmentColor;

void main()
{
    fragmentColor = vec4(color.r/0xFFFF, color.g/0xFFFF, color.b/0xFFFF, color.a/0xFFFF);
    gl_Position =  projection * view * model * vec4(pos, 1.0);
}

