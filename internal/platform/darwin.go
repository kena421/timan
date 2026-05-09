package platform

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void makeWindowTopmostAndFrameless(const char* title) {
    @autoreleasepool {
        NSString* nsTitle = [NSString stringWithUTF8String:title];
        dispatch_async(dispatch_get_main_queue(), ^{
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
                    break;
                }
            }
        });
    }
}
*/
import "C"
import "time"

func TweakWindow(title string) {
	// Give the window time to initialize
	go func() {
		time.Sleep(200 * time.Millisecond)
		C.makeWindowTopmostAndFrameless(C.CString(title))
	}()
}
