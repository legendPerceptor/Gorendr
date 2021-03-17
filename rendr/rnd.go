package rendr
/* The rndKernel stores everything about a reconstruction kernel. The kernel
   is non-zero only within [-support/2,support/2], for integer "support"
   which may be odd or even (but always positive). The kernels are set up at
   compile-time in such a way that each kernel knows its own derivative; the
   derivative of rndKernel *k is k->deriv. */
type evalFunc func(xx float64) float64
type applyFunc func(ww []float64, xa float64)
type rndKernel struct {
	Name, Desc string // identifying and descriptive string
	Support int // # samples needed for convolution
	Eval evalFunc // evaluate the kernel once
	Apply applyFunc // evaluate the kernel support times
	Deriv *rndKernel // derivative of this kernel; will point back to itself when kernel is zero
}

// kernels defined in kernel.go
//var rndKernelZero *rndKernel = &zero1Kernel
//var rndKernelCtmr *rndKernel = &ctmr

const RND_KERNEL_SUPPORT_MAX int = 6

// kernel parse
//func rndKernelParse(_kstr string) *rndKernel;

type rndType int
//rndType
const (
	rndTypeUnknown rndType=iota          // (0) no type known
	rndTypeUChar            // (1) uchar (8-bit unsigned)
	rndTypeShort            // (2) short (16-bit signed)
	rndTypefloat64             // (2) float64 (floating point)
)

type rndSpace int
//rndSpace
const (
	rndSpaceUnknown rndSpace=iota
	rndSpaceRGB
	rndSpaceHSV
	rndSpaceAlpha
)

type rndProbe int
//rndProbe
const (
	rndProbeUnknown rndProbe=iota /* (0) */
	rndProbePosView   /* (1) 3-vector: view-space position of probe
   						(for debugging ray-casting geometry) */
	rndProbePosWorld  /* (2) 3-vector: world-space position of probe
						(for debugging ray-casting geometry) */
	rndProbePosIndex  /* (3) 3-vector: index-space position of probe
   						(for debugging ray-casting geometry) */
	rndProbeInside    /* (4) scalar: 1 if kernel support entirely inside data,
   						0 if outside (copied from rndConvo->inside) */
	rndProbeValue     /* (5) scalar: value measured by convolution */
	rndProbeGradVec   /* (6) 3-vector: gradient (in world-space) measured by
   						convolution */
	rndProbeGradMag   /* (7) scalar: gradient magnitude */
	// in both of the following, the "a" is opacity *after* opacity correction
	rndProbeRgba      /* (8) 4-vector: RGBA from the transfer function: the
						   univariate LUT set by rndCtxTxfSet(). Also, if
						   rndCtxLevoySet() is called, the opacity is also
						   multiplied by the Levoy opacity functions. The RGB
						   part of this is the Blinn-Phong "material color" */
	rndProbeRgbaLit   /* (9) 4-vector: rgba from transfer function (above),
	   with RGB shaded according to Blinn-Phong, and then
	   multiplied by the lerp'd depth cueing color */
)

type rndBlend int
//rndBlend
const (
	rndBlendUnknown rndBlend=iota/* (0) blending not known */
	rndBlendMax      /* (1) max of all values (in vectors, per-component) */
	rndBlendSum      /* (2) sum of all values (in vectors, per-component)
					   It is possible to do max and sum blending of
					   probes rndProbeRgba and rndProbeRgbaLit. */
	rndBlendOver     /* (3) "over" operator of RGB color and alpha. This
					   is the one blend that breaks the orthogonality of
					   blending and probing; rndBlendOver only makes sense
					   for rndProbeRgba and rndProbeRgbaLit */
)

/*
  struct rndCamera: A container for information about how a camera is viewing
  some part of world-space.  The look-from, look-at, and up vectors are in
  world-space coordinates (using an orthonormal basis).

  All the fields are named here the same as in FSV Section 5.3
*/
type rndCamera struct {
	Fr,  // look-from point
	At,  // look-at point
	Up [3]float64 // up
	Nc, Fc, 		/* near,far clip plane distances,
	   				relative to the look-at point */
	FOV float64   // vertical field-of-view, in degrees
	Size [2]uint /* # horz,vert samples of image plane; this
	   				determines aspect ratio ("ar") */
	Ortho int /* if non-zero, use orthographic instead of
			   perspective projection. This does not affect how
			   any variables below are calculated, but the
			   rndCamera is best place to store this info. */
	// "output" parameters set by rndCameraUpdate(), based on the above
	Ar,        // ((float64)size[0])/size[1]
	D,           // distance between fr and at
	Ncv, Fcv,    /* (positive) distances, from eye to near and far
	   clipping planes */
	Hght, Wdth float64  // height and width of visible image plane
	U,        // right-ward basis vector of view space
	V,        // up-ward basis vector of view space
	N [3]float64        // back-ward (into eye) basis vector of view space
	VtoW [16]float64    // view-to-world transorm
}

/*
  struct rndVolume: A container for the oriented volume data that is
  processed by rendr.  Orientation is defined in terms of a particular
  orthonormal world-space basis: {left, posterior, superior}, which has
  significance for biomedical scans of human (or human-ish) anatomy, and can
  be understood as some right-handed coordinate frame otherwise.

  Data is always scalar (we are not handling multiple values per voxel)
*/
type rndVolume struct {
	Content string       /* if non-NULL, a descriptive string of what is in
	   this volume, or how it was generated */
	Size [3]uint         /* # samples along fastest (size[0]) to slowest
	   (size[2]) axes */
	ItoW [16]float64     /* homogeneous coordinate mapping from index-space
						   (fast-to-slow ordering) to the
						   "left-posterior-superior" world-space */
	Dtype rndType       /* type of the voxel data */

	Data *float64
}

/*
  rndImage is a container for output images. Can have 1, 2, 3, 4, or 5 values
  per pixel. The last component contains timing information if
  rndCtx->timing. No orientation information is saved. Implemented in image.c
*/
type rndImage struct {
	Content string      /* if non-NULL, a descriptive string of what is
	   						in this image, or how it was generated */
	Channel uint        /* how many values are at each pixel;
	   						this is always the fastest axis */
	Size [2] uint         /* # of samples along faster (size[0]) and
	   slower (size[1]) spatial axes */
	Dtype rndType       /* type of the data */

	Data *float64
}

/*
  struct rndConvo: the buffers and state that are used to do convolution and
  to store its results. This info is not in the rndCtx so that rndConvoEval
  can be thread-safe.  See go.c and ray.c for the context of how the rndConvo
  is passed to functions that may need to do convolution, like your function
  for handling each sample along a ray.
*/
type rndConvo struct{
	Ipos [3]float64 // index-space position position at last rndConvoEval
	Value float64
	Gradient [3]float64 // world-space gradient
	Inside int      /* 1 if the last call to rndConvoEval was at a position
				   where all the required data samples (in the support of
				   the kernel) were available, or 0 if not. This supplies
				   the value for rndProbeInside (above). Compared to
				   previous projects, there is no effort to count how
				   many samples were outside (relative to previous
				   projects (i.e. inside is the new !outside) */
	Kern_values,
	Deriv_values *float64 // for convolution computation
}



type _txf struct {
	len uint			// length of rgba LUT
	vmin, vmax float64	// domain of the LUT is the interval [vmin,vmax]
	rgba *float64		// lookup table data, rgba on faster axis
	unitStep float64	/* the "unit" length to use when computing opacity
					   correction as a function of ray step size
					   S. Whatever formula you use for opacity
					   correction should probably involve S/unitStep
					   rather than S alone; no such correction is
					   needed when S == unitStep. unitStep is by
					   default set (in ctx.c and rendr_go.c) to 1 */
	alphaNear1 float64 /* this value, which should be <= 1.0, controls
					   early ray termination. If < 1: with rndProbeRgba
					   and rndProbeRgbaLit, finish ray if its opacity
					   exceeds alphaNear1. Or, if == 1: no early ray
					   termination. */
}

type _levoy struct {
	num uint	/* number of surfaces; logical length of "vra"
			   array.  0 means "no Levoy opacity functions in
			   use" */
	vra *float64 /* 2-D array of parameters; v,r,a on faster axis;
			   logically a 1-D array of 3-vectors, where each
			   element of the 3-vector is:
			   v: isovalue, ("f_v" in Levoy paper),
			   r: fuzzy isocontour thickness,
			   a: max opacity of isocontour ("a_v" in paper) */
}

type _light struct {
	num uint	 /* number of lights, logical length of each of the
				   arrays below */
	rgb *float64 /* light ii has color (rgb + 3*ii)[0,1,2]. Colors
	   			are only in RGB space, not HSV */
	dir *float64 /* light direction (i.e. the direction *towards* the
				   light); this is not necessarily a normalized
				   vector, and it can be either in view-space or
				   world-space */
	vsp *int	 /* light ii is in view-space if vsp[ii] if is
	   				non-zero, else the light is in world-space */
	xyz *float64 /* normalized light direction in world space; set by
				   rndCtxLightUpdate() */
}

type _lparam struct {
	ka, kd, ks, p float64 /* Blinn-Phong parameters */
	dcn, dcf [3]float64 /* depth-cueing: as the last step of computing
						   rgbaLit probe values, the RGB color values
						   are multiplied (per-channel) by a lerp (maybe
						   with V3_LERP) between dcn and dcf as the ray
						   sampling position varies between the near and
						   far clipping planes; Setting both dcn and dcf
						   to (1,1,1) effectively turns off
						   depth-cueing */
}

/*
  struct rndCtx: all the state (the "context") associated with doing volume
  rendering. See rendr_go.c to see how the context is created and set up to
  perform rendering.

  Note: rndCtxGradientNeed(const rndCtx *ctx) is a useful little
  function that looks at ctx->probe and ctx->blend (after they are set
  by rndCtxParmSet()) to see if the gradient (or gradient magnitude)
  needs to be computed by rndConvoEval().
*/
type rndCtx struct {
	Vol *rndVolume
	Kern *rndKernel
	SampleOnceStop int /* for debugging: ray does only one sample, on near
	   clipping plane, via rndRayStart(), then stops */
	PlaneSep float64 /* separation (always > 0) between sampling planes
					   through the view space volume: all rays' N-th
					   samples lie within a plane N*planeSep behind the
					   near clipping plane. Thus, in perspective
					   projection, the step size on any single ray will
					   be slightly larger than planeSep except at the
					   very center of the image plane. */
	Probe rndProbe /* what should we be probing or computing at each
	  				 sample along each ray */
	Blend rndBlend /* how the probes along the ray wil be reduced
	  				 down to a single per-ray value or vector */
	ThreadNum uint /* how many threads to render with; 0 for
					   non-threaded execution */
	OutsideValue float64  /* the value that rndRayFinish should store in
				   ray->result, if there weren't any samples on the
				   ray that could be blended (e.g. when probing
				   rndProbeValue but the ray never went "inside"
				   the volume). If rndProbeLength(ctx->probe) > 1,
				   this outsideValue will be copied into each
				   channel. This can be NaN. */
	Timing int /* if non-zero: record per-ray timing (in
	   milliseconds) rndRay->time, and in an extra
	   last channel of the output image */
	// camera and image specification, set by rndCtxCameraSet()
	Cam rndCamera
	/* For rndProbeRgba and rndProbeRgbaLit: univariate RGBA transfer
	   function, represented as a lookup table "rgba", and other variables
	   related to opacity. Note: there are no rndTxf structs in here. Rather,
	   as a separate pre-process, the LUT is generated by rndTxfLutGenerate()
	   (via calls to rndTxfEval()), and saved to a file, which is then read
	   back in as the "-lut" option to "rendr go", and then passed to
	   rndCtxTxfSet(), which also sets unitStep and alphaNear1. */
	Txf _txf
	/* Additional *optional* opacity function, which multiples the opacity
	   generated by the univariate txf above: Levoy's "Isovalue contour
	   surface" opacity functions (from 1988 "Display of Surfaces from Volume
	   Data" paper; the triangular shapes in value/gradmag space). Parameters
	   are set by rndCtxLevoySet(). The evaluation of the opacity functions
	   is done by rndTxfLevoy(). Whether to use this transfer function is
	   testable by levoy.num being non-zero */
	Levoy _levoy
	/* Specification of multiple directional lights. The num, rgb, dir, and
	   vsp fields are all set by rndCtxLightSet(); xyz is allocated by
	   rndCtxLightSet() but the final values are set by rndCtxLightUpdate()
	   (which is called in rndRender()).  Thus, by the time you have to
	   compute lighting in ray.c, the only fields that matter are: num, rgb,
	   and xyz */
	Light _light
	/* Blinn-Phong and depth-cueing lighting parameters that are specific to
	   rndProbeRgbaLit; set by rndCtxLparmSet(). If light.num (above) is zero,
	   the only lighting is ambient (controlled by ka) and depth cueing.  */
	Lparam _lparam
	// student definition
	WtoI [16]float64
	MT [9]float64

}