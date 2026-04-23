width  = 450
  height = 200

  # Create one shape to be shared between the SVG and Canvas
  shapes = [0...3].map (i) ->
    shape = seen.Shapes.sphere(i).scale(height * 0.4)
    seen.Colors.randomSurfaces2(shape)
    return shape

  # Create one scene for each shape
  scenes = shapes.map (shape) -> new seen.Scene
    fractionalPoints : true
    model            : seen.Models.default().add(shape)
    viewport         : seen.Viewports.center(width, height)

  # Create a render context for each SVG and Canvas
  contexts = []
  for scene, i in scenes
    for type in ['canvas', 'svg']
      contexts.push seen.Context("seen-#{type}-#{i}", scene).render()

  # Slowly rotate shapes
  new seen.Animator().onFrame((t, dt) ->
    for shape in shapes then shape.rotx(dt*3e-4).roty(dt*2e-4)
    for context in contexts then context.render()
  ).start()