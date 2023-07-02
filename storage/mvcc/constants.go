// ---

package mvcc

import "math"

// Key layout:
// ...
// Key_lock        -- (0)
// Key_verMax      -- (1)
// ...
// Key_ver+1       -- (2)
// Key_ver         -- (3)
// Key_ver-1       -- (4)
// ...
// Key_0           -- (5)
// NextKey_lock    -- (6)
// NextKey_verMax  -- (7)
// ...
// NextKey_ver+1   -- (8)
// NextKey_ver     -- (9)
// NextKey_ver-1   -- (10)
// ...
// NextKey_0       -- (11)
// ...
// EOF

const LockVer = math.MaxUint64
