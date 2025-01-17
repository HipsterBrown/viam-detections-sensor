package detections_sensor

import (
	"context"
	"errors"
	"fmt"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/vision"
	"go.viam.com/utils/rpc"
)

var (
	Detections       = resource.NewModel("hipsterbrown", "sensor", "detections")
	errUnimplemented = errors.New("unimplemented")
)

func init() {
	resource.RegisterComponent(sensor.API, Detections,
		resource.Registration[sensor.Sensor, *Config]{
			Constructor: newDetectionsSensorDetections,
		},
	)
}

type Config struct {
	Camera   string   `json:"Camera"`
	Detector string   `json:"detector"`
	Labels   []string `json:"labels,omitempty"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *Config) Validate(path string) ([]string, error) {
	var deps []string
	if len(cfg.Camera) == 0 {
		return nil, resource.NewConfigValidationFieldRequiredError(path, "camera")
	}

	if len(cfg.Detector) == 0 {
		return nil, resource.NewConfigValidationFieldRequiredError(path, "detector")
	}

	deps = append(deps, cfg.Camera, cfg.Detector)
	return deps, nil
}

type detectionsSensorDetections struct {
	name resource.Name

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()

	detector vision.Service

	resource.TriviallyReconfigurable
	resource.TriviallyCloseable
}

func newDetectionsSensorDetections(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (sensor.Sensor, error) {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	s := &detectionsSensorDetections{
		name:       rawConf.ResourceName(),
		logger:     logger,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}

	if err := s.Reconfigure(ctx, deps, rawConf); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *detectionsSensorDetections) Name() resource.Name {
	return s.name
}

func (s *detectionsSensorDetections) Reconfigure(ctx context.Context, deps resource.Dependencies, rawConf resource.Config) error {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return err
	}
	s.cfg = conf

	detector, err := vision.FromDependencies(deps, conf.Detector)
	if err != nil {
		return fmt.Errorf("no detector for detections  (%s): %w", conf.Detector, err)
	}
	s.detector = detector

	return nil
}

func (s *detectionsSensorDetections) NewClientFromConn(ctx context.Context, conn rpc.ClientConn, remoteName string, name resource.Name, logger logging.Logger) (sensor.Sensor, error) {
	panic("not implemented")
}

func (s *detectionsSensorDetections) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	readings := make(map[string]interface{})
	labels := s.cfg.Labels

	for _, label := range labels {
		readings[label] = int64(0)
	}

	detections, err := s.detector.DetectionsFromCamera(ctx, s.cfg.Camera, nil)
	if err != nil {
		s.logger.Info("Detections error: ", readings)
		return readings, err
	}

	for _, detection := range detections {
		label := detection.Label()
		if len(labels) > 0 && !contains(labels, label) {
			continue
		}

		reading, ok := readings[label]
		if !ok {
			readings[label] = 1
		} else {
			if count, ok := reading.(int64); ok {
				readings[label] = count + 1
			} else {
				s.logger.Info("Unable to update count", label, reading)
			}
		}
	}
	return readings, nil
}

func (s *detectionsSensorDetections) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, errUnimplemented
}

func (s *detectionsSensorDetections) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
