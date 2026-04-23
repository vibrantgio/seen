width           = 900
  height          = 500
  petalPercentage = 0.03

  model    = new seen.Model()
  material = new seen.Material(seen.Colors.rgb(0xFF, 0xAA, 0xAA), {specularExponent : 1})
  rand     = -> (2*Math.random() - 1)

  # Lights
  model.add seen.Lights.directional
    normal    : seen.P(-1, 1, 0.5).normalize()
    color     : seen.Colors.hex('#FFFF00')
    intensity : 0.008

  model.add seen.Lights.ambient
    intensity : 0.004

  # Load petals
  petalShapes = []
  petalLoads  = ['assets/petal0.obj', 'assets/petal1.obj'].map (url) -> 
    return $.get(url, {}).then (contents) ->
      shapeFactory = ->
        shape = seen.Shapes.obj(contents, false)
        for surface in shape.surfaces
          surface.cullBackfaces = false
          surface.fill          = material
        return shape
      petalShapes.push(shapeFactory)
  petalsReady = $.when.apply($, petalLoads)

  generatePetalShape = (rho0, rho1, speed) ->
    shape = petalShapes[Math.floor(Math.random()*petalShapes.length)]()
      .scale(3)
      .rotz(rand())
      .bake()
    seen.Util.defaults(shape, {rho0, rho1, speed})
    return shape

  generateRandomPetal = ->
    model
      .append().translate(-width*0.8, (1.5*Math.random()-0.5)*height/2, Math.random()*-500)
      .append().translate(0, rand()*50)
      .add(generatePetalShape(rand(), rand(), Math.random() + 1))

  # Create scene
  scene = new seen.Scene
    model    : model
    viewport : seen.Viewports.center(width, height)

  # Create render context from canvas
  context = seen.Context('seen-group', scene).render()

  # Create animator
  animator = context.animate()
    .onBefore((t, dt) ->
      # Update petal animation
      for group0 in model.children
        group1 = group0.children[0]
        shape  = group1.children[0]
        group0.translate dt * 2e-1 * shape.speed
        group1.rotx dt * 3e-3 * shape.rho0
        shape.rotx dt * 1e-3 * shape.rho1

        shape.rho0 += rand()*1e-1
        shape.rho1 += rand()*1e-1

      # Occassionally add a new petal
      if Math.random() < petalPercentage then generateRandomPetal()
    )
    .onAfter(->
      # Remove petals that have blown out of view
      toRemove = model.children.filter (group0)->
        group1 = group0.children[0]
        shape  = group1.children[0]
        renderModel = scene._renderModelCache[shape.surfaces[0].id]
        return renderModel.projected.barycenter.x > width
      model.remove(toRemove...)
      if toRemove.length > 0 then scene.flushCache()
    )

  # Start animation once petal shapes are loaded
  petalsReady.then -> animator.start()
