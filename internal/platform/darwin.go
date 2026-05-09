package platform

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

// makeWindowTopmostAndFrameless finds a window by its title and applies 
// advanced macOS windowing properties to transform it into a persistent HUD.
//
// Key transformations:
// 1. Removes all title bar and window decorations (Borderless).
// 2. Elevates window level to NSPopUpMenuWindowLevel (above almost all apps).
// 3. Prevents hiding when the app is deactivated or "Hide Others" is used.
// 4. Sets collection behavior to allow visibility on Full Screen apps and all Spaces.
int makeWindowTopmostAndFrameless(const char* title) {
    __block int found = 0;
    @autoreleasepool {
        NSString* nsTitle = [NSString stringWithUTF8String:title];
        // Ensure UI updates happen on the main thread
        dispatch_sync(dispatch_get_main_queue(), ^{
            NSArray* windows = [NSApp windows];
            for (NSWindow* window in windows) {
                if ([[window title] isEqualToString:nsTitle]) {
                    // Force frameless state
                    [window setStyleMask:NSWindowStyleMaskBorderless];
                    
                    // High-priority level (Level 101) to stay above most windows
                    [window setLevel:NSPopUpMenuWindowLevel]; 
                    
                    // Persistence settings
                    [window setHidesOnDeactivate:NO];
                    [window setCanHide:NO];
                    
                    // Visual styling
                    [window setBackgroundColor:[NSColor clearColor]];
                    [window setOpaque:NO];
                    [window setHasShadow:YES];
                    [window setMovableByWindowBackground:YES];
                    
                    // Space and Full Screen behavior
                    [window setCollectionBehavior:NSWindowCollectionBehaviorCanJoinAllSpaces | 
                                                 NSWindowCollectionBehaviorFullScreenAuxiliary | 
                                                 NSWindowCollectionBehaviorIgnoresCycle | 
                                                 NSWindowCollectionBehaviorFullScreenDisallowsTiling];
                    found = 1;
                    break;
                }
            }
        });
    }
    return found;
}

// setWindowSharing toggles the window's visibility to screen capture APIs.
// Setting allow=0 uses NSWindowSharingNone to hide from Zoom, Slack, etc.
void setWindowSharing(const char* title, int allow) {
    @autoreleasepool {
        NSString* nsTitle = [NSString stringWithUTF8String:title];
        dispatch_async(dispatch_get_main_queue(), ^{
            NSArray* windows = [NSApp windows];
            for (NSWindow* window in windows) {
                if ([[window title] isEqualToString:nsTitle]) {
                    if (allow) {
                        [window setSharingType:NSWindowSharingReadOnly];
                    } else {
                        // NSWindowSharingNone makes the window invisible to screen sharing
                        [window setSharingType:NSWindowSharingNone];
                    }
                    break;
                }
            }
        });
    }
}

// setAccessoryMode changes the application's activation policy to 'Accessory'.
// This removes the app from the Dock and allows it to appear over Full Screen spaces.
void setAccessoryMode() {
    dispatch_async(dispatch_get_main_queue(), ^{
        [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
    });
}
*/
import "C"
import "time"

// TweakWindow manages the lifecycle of the window's HUD transformation.
// It uses an initial aggressive poll to eliminate startup flicker and a 
// persistent heartbeat to ensure settings stick across Space transitions.
func TweakWindow(title string) {
	// Set Accessory Mode immediately to allow overlays on Full Screen apps
	C.setAccessoryMode()

	go func() {
		cTitle := C.CString(title)
		
		// 1. Initial aggressive poll (1s)
		// Re-applies settings every 5ms to "win" any race conditions with the UI framework boot.
		for i := 0; i < 200; i++ {
			C.makeWindowTopmostAndFrameless(cTitle)
			time.Sleep(5 * time.Millisecond)
		}
		
		// 2. Persistent maintenance (every 1s)
		// Re-applies settings to ensure the HUD follows the user to new Full Screen apps/Spaces.
		for {
			C.makeWindowTopmostAndFrameless(cTitle)
			time.Sleep(time.Second)
		}
	}()
}

// SetPrivacyMode toggles whether the window is visible to screen capture/sharing APIs.
// When enabled=true, the window level is set to NSWindowSharingNone.
func SetPrivacyMode(title string, enabled bool) {
	cTitle := C.CString(title)
	allow := 1
	if enabled {
		allow = 0 // sharing none
	}
	C.setWindowSharing(cTitle, C.int(allow))
}
