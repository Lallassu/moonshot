#version 330 core
layout(location = 0) in vec4 xyzs; 
layout(location = 1) in vec4 color; 

uniform mat4 model, view, projection;

out vec4 particlecolor;
out vec2 pos;
out float size;           

void main()            
{                      
    gl_PointSize = xyzs.z;
    size = xyzs.w;
    
    pos = xyzs.xy;

    gl_Position = projection * view * vec4(xyzs.xyz, 1.0);
    particlecolor = color;
}


