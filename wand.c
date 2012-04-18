// Last updated 2008/11/04 10:53

// A simple example equivalent to
// convert logo: logo.jpg

#include <stdio.h>
#include <wand/magick_wand.h>

int main(int argc, char **argv) {
  MagickWand *mw = NULL;

  MagickWandGenesis();

  /* Create a wand */
  mw = NewMagickWand();

  /* Read the input image */
  MagickReadImage(mw,"example.jpg");

  /* write it */
  MagickWriteImage(mw,"example-c.png");

  /* Tidy up */
  if(mw) mw = DestroyMagickWand(mw);

  MagickWandTerminus();

  return 0;
}
