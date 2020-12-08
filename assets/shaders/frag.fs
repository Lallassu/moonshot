#version 330 core

in vec4 fragmentColor;

layout(location = 0) out vec4 color;

void main(void)
{
    color = fragmentColor;
}

