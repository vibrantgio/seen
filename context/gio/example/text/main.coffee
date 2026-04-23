width  = 900
height = 500

# Generate some random data points
data = [0...10].map -> Math.random() * 80.0 + 20.0

# Create scene model
model = seen.Models.default()

# Draw bars for data
for d, i in data
model.add(seen.Shapes.unitcube()
    .scale(20, d, 20)
    .translate(i * 30)
    .fill('#0088FF')
)

# Draw text above bars
for d, i in data
model.add(seen.Shapes.text(d.toFixed(1), {font : '10px Roboto', cullBackfaces : false, anchor : 'center'})
    .translate(i * 30 + 10, d + 10, 10)
    .fill('#000000')
)

# Create scene
scene = new seen.Scene
model    : model.translate(-150, -50, 0).scale(2)
viewport : seen.Viewports.center(width, height)

# Create render context from canvas
context = seen.Context('seen-canvas', scene).render()

# Enable drag-to-rotate on the canvas
dragger = new seen.Drag('seen-canvas', {inertia : true})
dragger.on('drag.rotate', (e) ->
xform = seen.Quaternion.xyToTransform(e.offsetRelative...)
model.transform(xform)
context.render()
)
