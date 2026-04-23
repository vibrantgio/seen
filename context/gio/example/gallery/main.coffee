  width  = 900
  height = 500

  # Create empty scene and render context
  scene = new seen.Scene
    model    : seen.Models.default()
    viewport : seen.Viewports.center(width, height)

  context = seen.Context('seen-canvas', scene).render()

  # Create pre-scaled group to add the shape to
  group = scene.model.append().scale(100)
  shape = null

  # Enable drag-to-rotate
  dragger = new seen.Drag(document.getElementById('seen-canvas'), {inertia : true})
  dragger.on('drag.rotate', (e) ->
    xform = seen.Quaternion.xyToTransform(e.offsetRelative...)
    group.transform(xform)
    context.render()
  )

  # Factories for generating sample shapes
  shapeFactory =
    sphere      : -> seen.Shapes.sphere(3).scale(2)
    tetrahedron : -> seen.Shapes.tetrahedron()
    cube        : -> seen.Shapes.cube()
    patch       : -> seen.Shapes.patch(4, 4).translate(-1, -1).scale(0.5)
    pipe        : -> seen.Shapes.pipe(seen.P(-1, -1, -1), seen.P(1, 1, 1), 0.2).scale(0.5)

  # =============================================
  #
  # Methods for updating scene from UI
  #
  # =============================================

  updateHue = ->
    hue = parseFloat($('#hue-slider').slider('value')) / 100.0
    for surf in shape.surfaces
      surf.fillMaterial.color = seen.Colors.hsl(hue, 0.5, 0.5)
      surf.fillMaterial.specularColor = surf.fillMaterial.color
      surf.dirty = true
    context.render()

  updateShinyness = ->
    shinyness = parseFloat($('#shinyness-slider').slider('value'))
    for surf in shape.surfaces
      surf.fillMaterial.specularExponent = shinyness
      surf.dirty = true
    context.render()

  updateMetallic = ->
    metallic = $('input[name=metallic-radio]:checked').val() is 'yes'
    for surf in shape.surfaces
      surf.fillMaterial.metallic = metallic
      surf.dirty = true
    context.render()

  updateShape =  ->
    shapeType = $('input[name=shape-radio]:checked').val()
    group.children = [shape = shapeFactory[shapeType]()]
    updateHue()
    updateShinyness()
    updateMetallic()
    context.render()

  updateLights = ->
    states = $('#light-toggles input').map (i, el) -> $(el).is(':checked')
    for state, i in states
      scene.model.lights[i].enabled = state
    context.render()

  updateLightingModel = ->
    shaderType = $('input[name=model-radio]:checked').val()
    scene.shader = seen.Shaders[shaderType]()
    context.render()

  # UI Events
  $(document).ready ->
    $('#shape-radios').buttonset().click(updateShape)
    $('#hue-slider').slider(slide : updateHue)
    $('#shinyness-slider').slider().slider(
      slide : updateShinyness
      value : 15
      min   : 1
      max   : 60
      step  : 1
    )
    $('#metallic-radios').buttonset().click(updateMetallic)
    $('#light-toggles').buttonset().click(updateLights)
    $('#lighting-model-radios').buttonset().click(updateLightingModel)

    # Initialize scene with UI values
    updateShape()
