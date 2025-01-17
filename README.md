# `detections-sensor` module

This [module](https://docs.viam.com/registry/modular-resources/) implements the [`rdk:component:sensor` API](https://docs.viam.com/appendix/apis/components/sensor/) for capturing the number of detections for each label from a vision service.

## Requirements

This module assumes you have an existing vision service, like a [mlmodel](https://docs.viam.com/operate/reference/services/vision/mlmodel/) or [YOLOv8](https://github.com/viam-labs/YOLOv8).

## Configure your detections sensor component

Navigate to the [**CONFIGURE** tab](https://docs.viam.com/configure/) of your [machine](https://docs.viam.com/fleet/machines/) in the [Viam app](https://app.viam.com/).
[Add `hipsterbrown:detections-sensor` to your machine](https://docs.viam.com/configure/#components).

### Attributes

The following attributes are available for `hipsterbrown:sensor:detections` sensor component:

| Name    | Type   | Required?    | Default | Description |
| ------- | ------ | ------------ | ------- | ----------- |
| `camera` | string | Required | N/A  | Name of the camera component to use with the vision service  |
| `detector` | string | Required | N/A  | Name of the vision service to use for getting detections  |
| `labels` | array of objects * | Optional     | [] | A list of specific labels to capture, if left empty all labels will be captured |

### Example configuration

```json
{
    "camera": "camera-1",
    "detector": "vision-1",
    "labels": ["person"]
}
```
