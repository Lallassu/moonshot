#version 330 core

in vec4 particlecolor;
in float size;
in vec2 pos;
layout(location = 0) out vec4 color;

void main() 
{
    color = particlecolor/0xFFFF;
}

