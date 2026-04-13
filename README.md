Rotating planet with a shader in Ebiten

![planet](./screenshot.jpg)

This is a simple example of how to use a shader in Ebiten to create a rotating planet effect. The shader takes as input a texture of a planet, which must be in equirectangular format (a rectangular image that can be mapped onto a sphere). Here is the earth image:

![earth](./earth.jpg)

For each pixel that the shader renders it does the following things:
1. Convert pixel coordinate to normalized texture coordinate in the range [-1, 1]
2. Convert the normalized texture coordinate to a 3D point on the unit sphere
3. Rotate the 3D point around an axis (a 3D vector) by an angle (in radians) that changes over time
4. Convert the rotated 3D point back to a normalized texture coordinate
5. Convert the normalized texture coordinate to a pixel coordinate in the input texture

And for good measure, apply a simple lighting effect to the pixel value. The planet rotates because in step 3 the angle changes over time, based on the number of ticks the engine has executed for.

Clouds are also generated on top of the planet, which uses the same exact renderer as the planet but with a smaller angle to make the clouds rotate slower than the planet, and also with a blend setting to make the clouds semi-transparent. The cloud texture is created by creating a new image the same size as the planet and copying small cloud images onto it at random positions.
