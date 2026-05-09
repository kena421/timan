package platform

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

int makeWindowTopmostAndFrameless(const char* title) {
    __block int found = 0;
    @autoreleasepool {
        NSString* nsTitle = [NSString stringWithUTF8String:title];
        // We use dispatch_sync to know immediately if we found the window
        dispatch_sync(dispatch_get_main_queue(), ^{
            NSArray* windows = [NSApp windows];
            for (NSWindow* window in windows) {
                if ([[window title] isEqualToString:nsTitle]) {
                    [window setStyleMask:NSWindowStyleMaskBorderless];
                    [window setLevel:NSStatusWindowLevel]; 
                    [window setBackgroundColor:[NSColor clearColor]];
                    [window setOpaque:NO];
                    [window setHasShadow:YES];
                    [window setMovableByWindowBackground:YES];
                    [window setCollectionBehavior:NSWindowCollectionBehaviorCanJoinAllSpaces | NSWindowCollectionBehaviorFullScreenAuxiliary];
                    found = 1;
                    break;
                }
            }
        });
    }
    return found;
}

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
                        [window setSharingType:NSWindowSharingNone];
                    }
                    break;
                }
            }
        });
    }
}
*/
import "C"
import "time"

// TweakWindow aggressively polls for the window by title to apply borderless styling 
// as soon as the window is managed by the OS, minimizing or eliminating title bar flicker.
func TweakWindow(title string) {
	go func() {
		cTitle := C.CString(title)
		// Poll every 5ms for up to 1 second
		for i := 0; i < 200; i++ {
			if C.makeWindowTopmostAndFrameless(cTitle) == 1 {
				return
			}
			time.Sleep(5 * time.Millisecond)
		}
	}()
}

// SetPrivacyMode toggles whether the window is visible to screen capture/sharing.
func SetPrivacyMode(title string, enabled bool) {
	cTitle := C.CString(title)
	allow := 1
	if enabled {
		allow = 0 // sharing none
	}
	C.setWindowSharing(cTitle, C.int(allow))
}
