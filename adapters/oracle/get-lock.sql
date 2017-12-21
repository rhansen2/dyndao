-- These must be loaded into your Oracle instance if you are going
-- to use locking.

-- They were taken from https://jeffkemponoracle.com/2005/10/user-named-locks-with-dbms_lock/

-- internal function to get a lock handle
-- (private for use by request_lock and release_lock)
CREATE OR REPLACE FUNCTION get_handle (lock_name IN VARCHAR2) RETURN VARCHAR2 IS
  PRAGMA AUTONOMOUS_TRANSACTION;
  lock_handle VARCHAR2(128);
BEGIN
  DBMS_LOCK.ALLOCATE_UNIQUE (
    lockname => lock_name,
    lockhandle => lock_handle,
    expiration_secs => 864000); -- 10 days
  RETURN lock_handle;
END get_handle;

CREATE OR REPLACE PROCEDURE request_lock (lock_name IN VARCHAR2) IS
  lock_status NUMBER;
BEGIN
  lock_status := DBMS_LOCK.REQUEST(
    lockhandle => get_handle(lock_name),
    lockmode => DBMS_LOCK.X_MODE, -- eXclusive
    timeout => DBMS_LOCK.MAXWAIT, -- wait forever
    release_on_commit => FALSE);
  CASE lock_status
    WHEN 0 THEN NULL;
    WHEN 2 THEN RAISE_APPLICATION_ERROR(-20000,'deadlock detected');
    WHEN 4 THEN RAISE_APPLICATION_ERROR(-20000,'lock already obtained');
    ELSE RAISE_APPLICATION_ERROR(-20000,'request lock failed: ' || lock_status);
  END CASE;
END request_lock;

-- wrapper to release a lock
CREATE OR REPLACE PROCEDURE release_lock (lock_name IN VARCHAR2) IS
  lock_status NUMBER;
BEGIN
  lock_status := DBMS_LOCK.RELEASE(
    lockhandle => get_handle(lock_name));
  IF lock_status > 0 THEN
    RAISE_APPLICATION_ERROR(-20000,'release lock failed: ' || lock_status);
  END IF;
END release_lock;

