# SDL Audio Support Implementation Plan

## Executive Summary
Add SDL audio support to RetroGoLib following the same CGO-free, modular design pattern established by the GUI package. The implementation will provide a clean interface for audio playback that can be used by emulators and retro console tooling.

## 1. Architecture Overview

### 1.1 Package Structure
```
audio/
├── audio.go          # Core interfaces and types
├── audio_test.go     # Unit tests for core functionality
├── doc.go            # Package documentation
└── sdl/
    ├── defines.go    # SDL audio constants and structures
    ├── imports.go    # Function imports and library loading
    ├── sdl.go        # Main SDL audio implementation
    ├── sdl_test.go   # SDL implementation tests
    ├── setup.go      # Unix/Linux setup
    └── setup_windows.go # Windows-specific setup
```

### 1.2 Design Principles
- **CGO-Free**: Use purego for dynamic library loading (same as GUI)
- **Modular**: Separate interface from implementation
- **Thread-Safe**: Handle audio callbacks and synchronization
- **Minimal Dependencies**: Only purego external dependency
- **Backend Agnostic**: Interface allows multiple audio backends

## 2. Core Interfaces

### 2.1 audio/audio.go
```go
package audio

// Format contains audio format specifications
type Format struct {
    SampleRate int     // e.g., 44100, 48000
    Channels   int     // 1 (mono) or 2 (stereo)
    Samples    int     // Buffer size in samples (e.g., 512, 1024)
    Format     uint16  // Sample format (e.g., int16, float32)
}

// Backend interface for audio emulation sources
type Backend interface {
    // AudioFormat returns the desired audio format
    AudioFormat() Format
    
    // AudioCallback fills the buffer with audio samples
    // Called by the audio system when more audio data is needed
    AudioCallback(buffer []byte)
    
    // AudioPaused returns whether audio should be paused
    AudioPaused() bool
}

// Initializer defines setup function for audio renderer
type Initializer func(backend Backend) (audioStart func() error, audioStop func(), err error)

// Setup will be set by the chosen audio renderer
var Setup Initializer
```

## 3. SDL Audio Implementation

### 3.1 audio/sdl/defines.go
```go
package sdl

// SDL Audio Format flags
const (
    AUDIO_U8     = 0x0008  // Unsigned 8-bit samples
    AUDIO_S8     = 0x8008  // Signed 8-bit samples
    AUDIO_U16LSB = 0x0010  // Unsigned 16-bit samples (little-endian)
    AUDIO_S16LSB = 0x8010  // Signed 16-bit samples (little-endian)
    AUDIO_U16MSB = 0x1010  // Unsigned 16-bit samples (big-endian)
    AUDIO_S16MSB = 0x9010  // Signed 16-bit samples (big-endian)
    AUDIO_S32LSB = 0x8020  // Signed 32-bit samples (little-endian)
    AUDIO_S32MSB = 0x9020  // Signed 32-bit samples (big-endian)
    AUDIO_F32LSB = 0x8120  // 32-bit floating point samples (little-endian)
    AUDIO_F32MSB = 0x9120  // 32-bit floating point samples (big-endian)
    
    // Platform-specific defaults
    AUDIO_U16SYS = AUDIO_U16LSB
    AUDIO_S16SYS = AUDIO_S16LSB
    AUDIO_S32SYS = AUDIO_S32LSB
    AUDIO_F32SYS = AUDIO_F32LSB
)

// SDL_AudioSpec structure
type AudioSpec struct {
    Freq     int32    // DSP frequency (samples per second)
    Format   uint16   // Audio data format
    Channels uint8    // Number of channels: 1 mono, 2 stereo
    Silence  uint8    // Audio buffer silence value
    Samples  uint16   // Audio buffer size in samples
    Padding  uint16   // Necessary for alignment
    Size     uint32   // Audio buffer size in bytes
    Callback uintptr  // Callback function pointer
    Userdata uintptr  // User data pointer
}

// SDL_AudioDeviceID type
type AudioDeviceID uint32

// Audio status
const (
    SDL_AUDIO_STOPPED = 0
    SDL_AUDIO_PLAYING = 1
    SDL_AUDIO_PAUSED  = 2
)
```

### 3.2 audio/sdl/imports.go
```go
package sdl

import (
    "fmt"
    "runtime"
    "github.com/ebitengine/purego"
)

var (
    // Core audio functions
    OpenAudioDevice func(device string, iscapture int, desired *AudioSpec, obtained *AudioSpec, allowedChanges int) AudioDeviceID
    CloseAudioDevice func(dev AudioDeviceID)
    PauseAudioDevice func(dev AudioDeviceID, pauseOn int)
    QueueAudio func(dev AudioDeviceID, data []byte, len uint32) int
    GetQueuedAudioSize func(dev AudioDeviceID) uint32
    ClearQueuedAudio func(dev AudioDeviceID)
    
    // Callback-based functions (if needed)
    LockAudioDevice func(dev AudioDeviceID)
    UnlockAudioDevice func(dev AudioDeviceID)
    
    // General SDL functions (shared with GUI)
    Init func(flags uint32) int
    GetError func() string
    Quit func()
)

var audioImports = map[string]any{
    "SDL_OpenAudioDevice":    &OpenAudioDevice,
    "SDL_CloseAudioDevice":   &CloseAudioDevice,
    "SDL_PauseAudioDevice":   &PauseAudioDevice,
    "SDL_QueueAudio":         &QueueAudio,
    "SDL_GetQueuedAudioSize": &GetQueuedAudioSize,
    "SDL_ClearQueuedAudio":   &ClearQueuedAudio,
    "SDL_LockAudioDevice":    &LockAudioDevice,
    "SDL_UnlockAudioDevice":  &UnlockAudioDevice,
}

// setupLibrary loads SDL library and registers functions
func setupLibrary() error {
    // Similar to GUI implementation
    // Check if SDL is already initialized by GUI
    // Load library and register audio-specific functions
}
```

### 3.3 audio/sdl/sdl.go
```go
package sdl

import (
    "fmt"
    "github.com/retroenv/retrogolib/audio"
)

// audioDevice holds the SDL audio device state
type audioDevice struct {
    deviceID AudioDeviceID
    backend  audio.Backend
    buffer   []byte
    running  bool
}

// Setup initializes SDL audio and returns control functions
func Setup(backend audio.Backend) (audioStart func() error, audioStop func(), err error) {
    if err := setupLibrary(); err != nil {
        return nil, nil, fmt.Errorf("setting up SDL audio library: %w", err)
    }
    
    // Initialize SDL audio subsystem if not already initialized
    if ret := Init(SDL_INIT_AUDIO); ret != 0 {
        return nil, nil, fmt.Errorf("initializing SDL audio: %s", GetError())
    }
    
    format := backend.AudioFormat()
    device := &audioDevice{
        backend: backend,
        buffer:  make([]byte, calculateBufferSize(format)),
    }
    
    // Configure audio specification
    desired := AudioSpec{
        Freq:     int32(format.SampleRate),
        Format:   format.Format,
        Channels: uint8(format.Channels),
        Samples:  uint16(format.Samples),
        Callback: 0, // Use queue-based audio (simpler than callbacks)
    }
    
    var obtained AudioSpec
    device.deviceID = OpenAudioDevice("", 0, &desired, &obtained, 0)
    if device.deviceID == 0 {
        return nil, nil, fmt.Errorf("opening audio device: %s", GetError())
    }
    
    start := func() error {
        return device.start()
    }
    
    stop := func() {
        device.stop()
    }
    
    return start, stop, nil
}

func (d *audioDevice) start() error {
    if d.running {
        return nil
    }
    
    d.running = true
    PauseAudioDevice(d.deviceID, 0) // Unpause
    
    // Start audio pumping goroutine
    go d.audioLoop()
    
    return nil
}

func (d *audioDevice) stop() {
    if !d.running {
        return
    }
    
    d.running = false
    ClearQueuedAudio(d.deviceID)
    CloseAudioDevice(d.deviceID)
}

func (d *audioDevice) audioLoop() {
    for d.running {
        // Check if we need to queue more audio
        queued := GetQueuedAudioSize(d.deviceID)
        if queued < uint32(len(d.buffer)*2) { // Keep 2 buffers queued
            // Get audio data from backend
            d.backend.AudioCallback(d.buffer)
            
            // Queue the audio data
            if ret := QueueAudio(d.deviceID, d.buffer, uint32(len(d.buffer))); ret != 0 {
                // Handle error
                continue
            }
        }
        
        // Check if audio should be paused
        if d.backend.AudioPaused() {
            PauseAudioDevice(d.deviceID, 1)
        } else {
            PauseAudioDevice(d.deviceID, 0)
        }
        
        // Small sleep to prevent busy-waiting
        time.Sleep(10 * time.Millisecond)
    }
}

func calculateBufferSize(format audio.Format) int {
    bytesPerSample := 2 // Default to 16-bit
    if format.Format == AUDIO_S32SYS || format.Format == AUDIO_F32SYS {
        bytesPerSample = 4
    } else if format.Format == AUDIO_U8 || format.Format == AUDIO_S8 {
        bytesPerSample = 1
    }
    
    return format.Samples * format.Channels * bytesPerSample
}
```

## 4. Integration Points

### 4.1 Emulator Integration Example
```go
// In an emulator implementation
type Emulator struct {
    // ... other fields
    audioBuffer []int16
    audioMutex  sync.Mutex
}

func (e *Emulator) AudioFormat() audio.Format {
    return audio.Format{
        SampleRate: 44100,
        Channels:   2,
        Samples:    512,
        Format:     sdl.AUDIO_S16SYS,
    }
}

func (e *Emulator) AudioCallback(buffer []byte) {
    e.audioMutex.Lock()
    defer e.audioMutex.Unlock()
    
    // Convert emulator's audio buffer to byte slice
    // Fill buffer with audio samples
}

func (e *Emulator) AudioPaused() bool {
    return e.paused
}
```

### 4.2 Main Application Setup
```go
import (
    "github.com/retroenv/retrogolib/audio"
    _ "github.com/retroenv/retrogolib/audio/sdl" // Register SDL audio
    "github.com/retroenv/retrogolib/gui"
    _ "github.com/retroenv/retrogolib/gui/sdl"   // Register SDL GUI
)

func main() {
    emulator := NewEmulator()
    
    // Setup GUI
    guiRender, guiCleanup, err := gui.Setup(emulator)
    if err != nil {
        log.Fatal(err)
    }
    defer guiCleanup()
    
    // Setup Audio
    audioStart, audioStop, err := audio.Setup(emulator)
    if err != nil {
        log.Fatal(err)
    }
    defer audioStop()
    
    // Start audio
    if err := audioStart(); err != nil {
        log.Fatal(err)
    }
    
    // Main loop
    for {
        running, err := guiRender()
        if err != nil || !running {
            break
        }
    }
}
```

## 5. Implementation Phases

### Phase 1: Core Structure (Week 1)
- [ ] Create audio package structure
- [ ] Define core interfaces (Backend, Format, Initializer)
- [ ] Add package documentation
- [ ] Write initial unit tests

### Phase 2: SDL Bindings (Week 1-2)
- [ ] Define SDL audio constants and structures
- [ ] Implement purego function imports
- [ ] Add platform-specific setup (Unix/Windows)
- [ ] Test library loading

### Phase 3: Audio Implementation (Week 2)
- [ ] Implement queue-based audio playback
- [ ] Add audio loop goroutine
- [ ] Handle pause/resume functionality
- [ ] Implement error handling and recovery

### Phase 4: Testing & Integration (Week 3)
- [ ] Write comprehensive unit tests
- [ ] Create integration tests with mock backend
- [ ] Test with actual emulator
- [ ] Performance testing and optimization

### Phase 5: Documentation & Polish (Week 3-4)
- [ ] Complete API documentation
- [ ] Add usage examples
- [ ] Create troubleshooting guide
- [ ] Update main README

## 6. Testing Strategy

### 6.1 Unit Tests
- Test Format struct validation
- Test Backend interface compliance
- Test buffer size calculations
- Test SDL function registration

### 6.2 Integration Tests
- Test with mock audio backend
- Test pause/resume functionality
- Test audio format negotiation
- Test cleanup and resource management

### 6.3 Performance Tests
- Benchmark audio callback performance
- Test latency measurements
- Memory usage profiling
- CPU usage monitoring

## 7. Considerations

### 7.1 Thread Safety
- Audio callbacks run in separate goroutine
- Mutex protection for shared state
- Lock-free audio buffer if possible

### 7.2 Error Handling
- Graceful degradation if audio unavailable
- Clear error messages for debugging
- Recovery from temporary failures

### 7.3 Platform Compatibility
- Test on Linux, macOS, Windows
- Handle different SDL2 versions
- Support various audio formats

### 7.4 Performance
- Minimize allocations in audio path
- Use appropriate buffer sizes
- Consider ring buffer for audio data

## 8. Alternative Approaches Considered

### 8.1 Callback-Based Audio
- More complex due to CGO limitations
- Would require unsafe pointer handling
- Decided on queue-based for simplicity

### 8.2 Multiple Backend Support
- Could add OpenAL, PortAudio backends
- Current design allows future expansion
- SDL2 provides good cross-platform support

### 8.3 Direct Hardware Access
- Platform-specific implementations
- More complex maintenance
- SDL2 abstraction preferred

## 9. Success Metrics

- Clean API matching GUI package design
- Zero CGO dependencies
- < 50ms audio latency
- Support for common sample rates (22050, 44100, 48000)
- Thread-safe operation
- Comprehensive test coverage (>80%)
- Works with existing RetroGoLib emulators

## 10. Future Enhancements

- Audio effects (reverb, filters)
- Multi-device support
- Recording capabilities
- MIDI support
- Advanced format conversion
- Hardware acceleration support