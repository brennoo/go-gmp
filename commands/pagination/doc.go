// Package pagination provides a generic, stateful iterator over GMP list
// commands (get_tasks, get_results, get_assets, get_targets, get_tickets,
// get_port_lists, get_settings), fetching successive pages transparently as
// the caller advances.
//
// [github.com/brennoo/go-gmp/client.Client] exposes one iterator constructor
// per resource — Tasks, Results, Assets, Targets, Tickets, PortLists,
// Settings — each returning an [Iterator] specialized for that resource's
// item type via a type alias (for example [TaskIterator]). Each accepts a
// page size and the same [github.com/brennoo/go-gmp/commands/filtering.FilterArg]
// values as the equivalent one-shot Get* method:
//
//	taskIter := cli.Tasks(ctx, 50, filtering.WithStatus(filtering.StatusRunning))
//	defer taskIter.Close()
//
//	for taskIter.Next() {
//		task := taskIter.Current()
//		fmt.Println(task.Name, task.Status)
//	}
//	if err := taskIter.Err(); err != nil {
//		log.Fatal(err)
//	}
//
// Next fetches a new page only when the current one is exhausted, so a caller
// that stops early (via break, or simply not calling Next again) issues no
// more requests than pages consumed. Err reports the error kind that stopped
// iteration — the same four kinds [github.com/brennoo/go-gmp/commands]
// returns from a one-shot call. Close releases the iterator's internal
// buffers; it does not close the underlying connection.
//
// [PaginationOptions] and [BuildPaginationFilter] are the lower-level pieces
// iterators are built from, useful when implementing a custom iterator for a
// resource this package does not already cover.
package pagination
