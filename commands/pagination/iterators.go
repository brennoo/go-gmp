//nolint:dupl // Similar patterns across iterators are intentional

package pagination

import (
	"github.com/brennoo/go-gmp/commands"
)

// TaskIterator implementation.

//nolint:dupl // Similar pattern across iterators is intentional
func (it *TaskIterator) loadPage() error {
	opts := PaginationOptions{
		Page:     it.Page,
		PageSize: it.Opts.PageSize,
	}
	filter := BuildPaginationFilter(opts, it.Filters...)

	cmd := &commands.GetTasks{
		Filter: filter,
	}

	resp, err := it.Client.GetTasksRaw(cmd)
	if err != nil {
		return err
	}

	// Convert tasks to typed slice
	it.current = make([]*commands.GetTasksResponseTask, len(resp.Task))
	for i, task := range resp.Task {
		it.current[i] = &task
	}

	// Update total if available (using TaskCount for now)
	if resp.TaskCount.Filtered > 0 {
		it.total = resp.TaskCount.Filtered
	}

	// Check if there are more pages
	it.HasMoreData = len(resp.Task) == it.Opts.PageSize

	return nil
}

// ResultIterator implementation.
func (it *ResultIterator) loadPage() error {
	opts := PaginationOptions{
		Page:     it.Page,
		PageSize: it.Opts.PageSize,
	}
	filter := BuildPaginationFilter(opts, it.Filters...)

	cmd := &commands.GetResults{
		Filter: filter,
	}

	resp, err := it.Client.GetResultsRaw(cmd)
	if err != nil {
		return err
	}

	// Convert results to typed slice
	it.current = make([]*commands.Result, len(resp.Results))
	for i, result := range resp.Results {
		it.current[i] = &result
	}

	// Update total if available
	if resp.Count != nil && resp.Count.Filtered > 0 {
		it.total = resp.Count.Filtered
	}

	// Check if there are more pages
	it.HasMoreData = len(resp.Results) == it.Opts.PageSize

	return nil
}

// AssetIterator implementation.
func (it *AssetIterator) loadPage() error {
	opts := PaginationOptions{
		Page:     it.Page,
		PageSize: it.Opts.PageSize,
	}
	filter := BuildPaginationFilter(opts, it.Filters...)

	cmd := &commands.GetAssets{
		Filter: filter,
	}

	resp, err := it.Client.GetAssetsRaw(cmd)
	if err != nil {
		return err
	}

	// Convert assets to typed slice
	it.current = make([]*commands.Asset, len(resp.Assets))
	for i, asset := range resp.Assets {
		it.current[i] = &asset
	}

	// Update total if available
	if resp.AssetCount != nil && resp.AssetCount.Filtered > 0 {
		it.total = resp.AssetCount.Filtered
	}

	// Check if there are more pages
	it.HasMoreData = len(resp.Assets) == it.Opts.PageSize

	return nil
}

// TargetIterator implementation.
func (it *TargetIterator) loadPage() error {
	opts := PaginationOptions{
		Page:     it.Page,
		PageSize: it.Opts.PageSize,
	}
	filter := BuildPaginationFilter(opts, it.Filters...)

	cmd := &commands.GetTargets{
		Filter: filter,
	}

	resp, err := it.Client.GetTargetsRaw(cmd)
	if err != nil {
		return err
	}

	// Convert targets to typed slice
	it.current = make([]*commands.Target, len(resp.Target))
	for i, target := range resp.Target {
		it.current[i] = &target
	}

	// Update total if available
	if resp.TargetCount != nil && resp.TargetCount.Filtered > 0 {
		it.total = resp.TargetCount.Filtered
	}

	// Check if there are more pages
	it.HasMoreData = len(resp.Target) == it.Opts.PageSize

	return nil
}

// TicketIterator implementation.
func (it *TicketIterator) loadPage() error {
	opts := PaginationOptions{
		Page:     it.Page,
		PageSize: it.Opts.PageSize,
	}
	filter := BuildPaginationFilter(opts, it.Filters...)

	cmd := &commands.GetTickets{
		Filter: filter,
	}

	resp, err := it.Client.GetTicketsRaw(cmd)
	if err != nil {
		return err
	}

	// Convert tickets to typed slice
	it.current = make([]*commands.Ticket, len(resp.Tickets))
	for i, ticket := range resp.Tickets {
		it.current[i] = &ticket
	}

	// Update total if available
	if resp.TicketCount != nil && resp.TicketCount.Filtered > 0 {
		it.total = resp.TicketCount.Filtered
	}

	// Check if there are more pages
	it.HasMoreData = len(resp.Tickets) == it.Opts.PageSize

	return nil
}

// PortListIterator implementation.
func (it *PortListIterator) loadPage() error {
	opts := PaginationOptions{
		Page:     it.Page,
		PageSize: it.Opts.PageSize,
	}
	filter := BuildPaginationFilter(opts, it.Filters...)

	cmd := &commands.GetPortLists{
		Filter: filter,
	}

	resp, err := it.Client.GetPortListsRaw(cmd)
	if err != nil {
		return err
	}

	// Convert port lists to typed slice
	it.current = make([]*commands.PortList, len(resp.PortLists))
	for i, portList := range resp.PortLists {
		it.current[i] = &portList
	}

	// Update total if available
	if resp.Count != nil && resp.Count.Filtered > 0 {
		it.total = resp.Count.Filtered
	}

	// Check if there are more pages
	it.HasMoreData = len(resp.PortLists) == it.Opts.PageSize

	return nil
}

// SettingsIterator implementation.
func (it *SettingsIterator) loadPage() error {
	opts := PaginationOptions{
		Page:     it.Page,
		PageSize: it.Opts.PageSize,
	}
	filter := BuildPaginationFilter(opts, it.Filters...)

	cmd := &commands.GetSettings{
		Filter: filter,
	}

	resp, err := it.Client.GetSettingsRaw(cmd)
	if err != nil {
		return err
	}

	// Convert settings to typed slice
	it.current = make([]*commands.Setting, len(resp.SettingsList))
	for i, setting := range resp.SettingsList {
		it.current[i] = &setting
	}

	// Update total if available
	if resp.SettingCount != nil && resp.SettingCount.Filtered > 0 {
		it.total = resp.SettingCount.Filtered
	}

	// Check if there are more pages
	it.HasMoreData = len(resp.SettingsList) == it.Opts.PageSize

	return nil
}
