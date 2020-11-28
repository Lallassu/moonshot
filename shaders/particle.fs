#version 330 core

in vec4 particlecolor;
in float size;
in vec2 pos;
out vec4 color;

void main() 
{
    //float d = distance(gl_FragCoord.xy, pos);
    color = particlecolor/0xFFFF;
    //if(d <= 0) {
    //    d = 0.1;
    //}
    //float a = d/size;
    //color.a = a;
}

