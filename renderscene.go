package seen

//----------------------------
// RenderScene
// implements RenderLayer
//----------------------------

type RenderScene struct {
	models []*RenderModel
}

func NewRenderScene() *RenderScene {
	return &RenderScene{}
}

func (s *RenderScene) Render(context RenderLayerContext) {
	for _, r := range s.models {
		r.Render(context)
	}
}

//----------------------------
// RenderModel
//----------------------------

type RenderModel struct {
	surface Surface
}

func (m *RenderModel) Render(context RenderLayerContext) {
	m.surface.Painter().Paint(m, context)
}

//----------------------------
// Surface
//----------------------------

type Surface interface {
	Painter() SurfacePainter
}

//----------------------------
// SurfacePainter
//----------------------------

type SurfacePainter interface {
	Paint(model *RenderModel, context RenderLayerContext)
}

//----------------------------
// SurfacePathPainter
//----------------------------

type SurfacePathPainter struct {
}

func NewSurfacePathPainter() SurfacePainter {
	return &SurfacePathPainter{}
}

func (p *SurfacePathPainter) Paint(model *RenderModel, context RenderLayerContext) {
	// painter := context.Path()
	// painter.Path(renderModel.projected.points)

	// if renderModel.fill?
	//   painter.Fill(
	//     fill           : if not renderModel.fill? then 'none' else renderModel.fill.hex()
	//     'fill-opacity' : if not renderModel.fill?.a? then 1.0 else (renderModel.fill.a / 255.0)
	//   )

	// if renderModel.stroke?
	//   painter.Draw(
	//     fill           : 'none'
	//     stroke         : if not renderModel.stroke? then 'none' else renderModel.stroke.hex()
	//     'stroke-width' : renderModel.surface['stroke-width'] ? 1
	//   )
}

//----------------------------
// SurfaceTextPainter
//----------------------------

type SurfaceTextPainter struct {
}

func NewSurfaceTextPainter() SurfacePainter {
	return &SurfaceTextPainter{}
}

func (p *SurfaceTextPainter) Paint(model *RenderModel, context RenderLayerContext) {
	// style = {
	//   fill          : if not renderModel.fill? then 'none' else renderModel.fill.hex()
	//   font          : renderModel.surface.font
	//   'text-anchor' : renderModel.surface.anchor ? 'middle'
	// }
	// xform = seen.Affine.solveForAffineTransform(renderModel.projected.points)
	// context.text().fillText(xform, renderModel.surface.text, style)
}
