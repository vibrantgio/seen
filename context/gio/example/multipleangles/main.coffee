
  dims = [
    [600, 300]
    [150, 150]
    [150, 150]
    [150, 150]
    [150, 150]
  ]

  # One model shared between scenes
  model  = seen.Models.default()

  scenes = [0...5].map (i) -> new seen.Scene
    fractionalPoints : true
    model            : model
    viewport         : seen.Viewports.center(dims[i]...)

  # Alter camera matrix for main viewport
  scenes[0].camera.scale(2)

  # Rotate mini viewports
  scenes[2].camera.rotx(Math.PI/2)
  scenes[3].camera.rotx(-Math.PI/2)
  scenes[4].camera.roty(-Math.PI/2)

  # Load the object file using jquery
  $.get 'assets/bunny-low.obj', {}, (contents) ->
    # Create shape from object file
    shape = seen.Shapes.obj(contents, false)
    shape.scale(8).translate(0,-30).rotx(Math.PI/4).roty(-Math.PI/4).rotz(-Math.PI/4)

    # Apply random colors to surfaces
    seen.Colors.randomSurfaces2(shape)

    # Update scene model
    model.add(shape)
    renderAll()

  # Create rendering contexts
  contexts = scenes.map (scene,i) -> seen.Context("seen-canvas-#{i}", scene)

  renderAll = ->
    context.render() for context in contexts

  # Enable drag-to-rotate on main viewport
  dragger = new seen.Drag('seen-canvas-0', {inertia : true})
  dragger.on('drag.rotate', (e) ->
    xform = seen.Quaternion.xyToTransform(e.offsetRelative...)
    model.transform(xform)
    renderAll()
  )