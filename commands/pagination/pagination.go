package pagination

import (
	"context"
	"fmt"

	"github.com/brennoo/go-gmp/commands"
)

// Client defines the minimal interface needed for pagination
type Client interface {
	GetTasks(cmd *commands.GetTasks) (*commands.GetTasksResponse, error)
	GetResults(cmd *commands.GetResults) (*commands.GetResultsResponse, error)
	GetAssets(cmd *commands.GetAssets) (*commands.GetAssetsResponse, error)
	GetTargets(cmd *commands.GetTargets) (*commands.GetTargetsResponse, error)
	GetTickets(cmd *commands.GetTickets) (*commands.GetTicketsResponse, error)
	GetPortLists(cmd *commands.GetPortLists) (*commands.GetPortListsResponse, error)
	GetSettings(cmd *commands.GetSettings) (*commands.GetSettingsResponse, error)
}

// PaginationOptions represents pagination configuration
type PaginationOptions struct {
	Page     int // 1-based page number
	PageSize int // Number of items per page
	MaxItems int // Maximum total items to fetch (0 = no limit)
}

// DefaultPaginationOptions returns sensible default pagination options
func DefaultPaginationOptions() PaginationOptions {
	return PaginationOptions{
		Page:     1,
		PageSize: 100,
		MaxItems: 0, // No limit
	}
}

// BuildPaginationFilter constructs the filter string for pagination
func BuildPaginationFilter(opts PaginationOptions, additionalFilters ...string) string {
	first := (opts.Page-1)*opts.PageSize + 1
	rows := opts.PageSize

	if opts.MaxItems > 0 && first > opts.MaxItems {
		return "rows=0" // No results if we're past the max
	}

	if opts.MaxItems > 0 && first+rows-1 > opts.MaxItems {
		rows = opts.MaxItems - first + 1
	}

	filter := fmt.Sprintf("first=%d rows=%d", first, rows)

	// Add additional filters
	for _, f := range additionalFilters {
		if f != "" {
			filter += " " + f
		}
	}

	return filter
}

// Iterator defines the generic iterator interface
type Iterator[T any] interface {
	Next() bool
	Current() T
	Err() error
	HasMore() bool
	Total() int
	Close()
}

// TaskIterator provides iteration over tasks with automatic pagination
type TaskIterator struct {
	Client      Client
	Ctx         context.Context
	Opts        PaginationOptions
	Filters     []string
	current     []*commands.GetTasksResponseTask
	index       int
	Page        int
	HasMoreData bool
	total       int
	err         error
}

// Next advances the iterator to the next item and returns true if successful
func (it *TaskIterator) Next() bool {
	if it.err != nil {
		return false
	}

	// Check context
	if it.Ctx != nil {
		select {
		case <-it.Ctx.Done():
			it.err = it.Ctx.Err()
			return false
		default:
		}
	}

	// Check if we need to load more data
	if it.index >= len(it.current) {
		if !it.HasMoreData {
			return false // No more data
		}

		// Load next page
		it.Opts.Page = it.Page
		it.err = it.loadPage()
		if it.err != nil {
			return false
		}

		if len(it.current) == 0 {
			it.HasMoreData = false
			return false // No more data
		}

		it.index = 0
		it.Page++
	}

	// Advance index for next call
	it.index++
	return true
}

// Current returns the current item (call after Next() returns true)
func (it *TaskIterator) Current() *commands.GetTasksResponseTask {
	if it.index == 0 || it.index > len(it.current) {
		return nil
	}
	return it.current[it.index-1]
}

// HasMore returns true if there are more items to iterate
func (it *TaskIterator) HasMore() bool {
	return it.HasMoreData || it.index < len(it.current)
}

// Err returns the last error encountered
func (it *TaskIterator) Err() error {
	return it.err
}

// Total returns the total number of items (if available)
func (it *TaskIterator) Total() int {
	return it.total
}

// Close cleans up the iterator
func (it *TaskIterator) Close() {
	it.HasMoreData = false
	it.current = nil
}

// ResultIterator provides iteration over results with automatic pagination
type ResultIterator struct {
	Client      Client
	Ctx         context.Context
	Opts        PaginationOptions
	Filters     []string
	current     []*commands.Result
	index       int
	Page        int
	HasMoreData bool
	total       int
	err         error
}

// Next advances the iterator to the next item and returns true if successful
func (it *ResultIterator) Next() bool {
	if it.err != nil {
		return false
	}

	// Check context
	if it.Ctx != nil {
		select {
		case <-it.Ctx.Done():
			it.err = it.Ctx.Err()
			return false
		default:
		}
	}

	// Check if we need to load more data
	if it.index >= len(it.current) {
		if !it.HasMoreData {
			return false // No more data
		}

		// Load next page
		it.Opts.Page = it.Page
		it.err = it.loadPage()
		if it.err != nil {
			return false
		}

		if len(it.current) == 0 {
			it.HasMoreData = false
			return false // No more data
		}

		it.index = 0
		it.Page++
	}

	// Advance index for next call
	it.index++
	return true
}

// Current returns the current item (call after Next() returns true)
func (it *ResultIterator) Current() *commands.Result {
	if it.index == 0 || it.index > len(it.current) {
		return nil
	}
	return it.current[it.index-1]
}

// HasMore returns true if there are more items to iterate
func (it *ResultIterator) HasMore() bool {
	return it.HasMoreData || it.index < len(it.current)
}

// Err returns the last error encountered
func (it *ResultIterator) Err() error {
	return it.err
}

// Total returns the total number of items (if available)
func (it *ResultIterator) Total() int {
	return it.total
}

// Close cleans up the iterator
func (it *ResultIterator) Close() {
	it.HasMoreData = false
	it.current = nil
}

// AssetIterator provides iteration over assets with automatic pagination
type AssetIterator struct {
	Client      Client
	Ctx         context.Context
	Opts        PaginationOptions
	Filters     []string
	current     []*commands.Asset
	index       int
	Page        int
	HasMoreData bool
	total       int
	err         error
}

// Next advances the iterator to the next item and returns true if successful
func (it *AssetIterator) Next() bool {
	if it.err != nil {
		return false
	}

	// Check context
	if it.Ctx != nil {
		select {
		case <-it.Ctx.Done():
			it.err = it.Ctx.Err()
			return false
		default:
		}
	}

	// Check if we need to load more data
	if it.index >= len(it.current) {
		if !it.HasMoreData {
			return false // No more data
		}

		// Load next page
		it.Opts.Page = it.Page
		it.err = it.loadPage()
		if it.err != nil {
			return false
		}

		if len(it.current) == 0 {
			it.HasMoreData = false
			return false // No more data
		}

		it.index = 0
		it.Page++
	}

	// Advance index for next call
	it.index++
	return true
}

// Current returns the current item (call after Next() returns true)
func (it *AssetIterator) Current() *commands.Asset {
	if it.index == 0 || it.index > len(it.current) {
		return nil
	}
	return it.current[it.index-1]
}

// HasMore returns true if there are more items to iterate
func (it *AssetIterator) HasMore() bool {
	return it.HasMoreData || it.index < len(it.current)
}

// Err returns the last error encountered
func (it *AssetIterator) Err() error {
	return it.err
}

// Total returns the total number of items (if available)
func (it *AssetIterator) Total() int {
	return it.total
}

// Close cleans up the iterator
func (it *AssetIterator) Close() {
	it.HasMoreData = false
	it.current = nil
}

// TargetIterator provides iteration over targets with automatic pagination
type TargetIterator struct {
	Client      Client
	Ctx         context.Context
	Opts        PaginationOptions
	Filters     []string
	current     []*commands.Target
	index       int
	Page        int
	HasMoreData bool
	total       int
	err         error
}

// Next advances the iterator to the next item and returns true if successful
func (it *TargetIterator) Next() bool {
	if it.err != nil {
		return false
	}

	// Check context
	if it.Ctx != nil {
		select {
		case <-it.Ctx.Done():
			it.err = it.Ctx.Err()
			return false
		default:
		}
	}

	// Check if we need to load more data
	if it.index >= len(it.current) {
		if !it.HasMoreData {
			return false // No more data
		}

		// Load next page
		it.Opts.Page = it.Page
		it.err = it.loadPage()
		if it.err != nil {
			return false
		}

		if len(it.current) == 0 {
			it.HasMoreData = false
			return false // No more data
		}

		it.index = 0
		it.Page++
	}

	// Advance index for next call
	it.index++
	return true
}

// Current returns the current item (call after Next() returns true)
func (it *TargetIterator) Current() *commands.Target {
	if it.index == 0 || it.index > len(it.current) {
		return nil
	}
	return it.current[it.index-1]
}

// HasMore returns true if there are more items to iterate
func (it *TargetIterator) HasMore() bool {
	return it.HasMoreData || it.index < len(it.current)
}

// Err returns the last error encountered
func (it *TargetIterator) Err() error {
	return it.err
}

// Total returns the total number of items (if available)
func (it *TargetIterator) Total() int {
	return it.total
}

// Close cleans up the iterator
func (it *TargetIterator) Close() {
	it.HasMoreData = false
	it.current = nil
}

// TicketIterator provides iteration over tickets with automatic pagination
type TicketIterator struct {
	Client      Client
	Ctx         context.Context
	Opts        PaginationOptions
	Filters     []string
	current     []*commands.Ticket
	index       int
	Page        int
	HasMoreData bool
	total       int
	err         error
}

// Next advances the iterator to the next item and returns true if successful
func (it *TicketIterator) Next() bool {
	if it.err != nil {
		return false
	}

	// Check context
	if it.Ctx != nil {
		select {
		case <-it.Ctx.Done():
			it.err = it.Ctx.Err()
			return false
		default:
		}
	}

	// Check if we need to load more data
	if it.index >= len(it.current) {
		if !it.HasMoreData {
			return false // No more data
		}

		// Load next page
		it.Opts.Page = it.Page
		it.err = it.loadPage()
		if it.err != nil {
			return false
		}

		if len(it.current) == 0 {
			it.HasMoreData = false
			return false // No more data
		}

		it.index = 0
		it.Page++
	}

	// Advance index for next call
	it.index++
	return true
}

// Current returns the current item (call after Next() returns true)
func (it *TicketIterator) Current() *commands.Ticket {
	if it.index == 0 || it.index > len(it.current) {
		return nil
	}
	return it.current[it.index-1]
}

// HasMore returns true if there are more items to iterate
func (it *TicketIterator) HasMore() bool {
	return it.HasMoreData || it.index < len(it.current)
}

// Err returns the last error encountered
func (it *TicketIterator) Err() error {
	return it.err
}

// Total returns the total number of items (if available)
func (it *TicketIterator) Total() int {
	return it.total
}

// Close cleans up the iterator
func (it *TicketIterator) Close() {
	it.HasMoreData = false
	it.current = nil
}

// PortListIterator provides iteration over port lists with automatic pagination
type PortListIterator struct {
	Client      Client
	Ctx         context.Context
	Opts        PaginationOptions
	Filters     []string
	current     []*commands.PortList
	index       int
	Page        int
	HasMoreData bool
	total       int
	err         error
}

// Next advances the iterator to the next item and returns true if successful
func (it *PortListIterator) Next() bool {
	if it.err != nil {
		return false
	}

	// Check context
	if it.Ctx != nil {
		select {
		case <-it.Ctx.Done():
			it.err = it.Ctx.Err()
			return false
		default:
		}
	}

	// Check if we need to load more data
	if it.index >= len(it.current) {
		if !it.HasMoreData {
			return false // No more data
		}

		// Load next page
		it.Opts.Page = it.Page
		it.err = it.loadPage()
		if it.err != nil {
			return false
		}

		if len(it.current) == 0 {
			it.HasMoreData = false
			return false // No more data
		}

		it.index = 0
		it.Page++
	}

	// Advance index for next call
	it.index++
	return true
}

// Current returns the current item (call after Next() returns true)
func (it *PortListIterator) Current() *commands.PortList {
	if it.index == 0 || it.index > len(it.current) {
		return nil
	}
	return it.current[it.index-1]
}

// HasMore returns true if there are more items to iterate
func (it *PortListIterator) HasMore() bool {
	return it.HasMoreData || it.index < len(it.current)
}

// Err returns the last error encountered
func (it *PortListIterator) Err() error {
	return it.err
}

// Total returns the total number of items (if available)
func (it *PortListIterator) Total() int {
	return it.total
}

// Close cleans up the iterator
func (it *PortListIterator) Close() {
	it.HasMoreData = false
	it.current = nil
}

// SettingsIterator provides iteration over settings with automatic pagination
type SettingsIterator struct {
	Client      Client
	Ctx         context.Context
	Opts        PaginationOptions
	Filters     []string
	current     []*commands.Setting
	index       int
	Page        int
	HasMoreData bool
	total       int
	err         error
}

// Next advances the iterator to the next item and returns true if successful
func (it *SettingsIterator) Next() bool {
	if it.err != nil {
		return false
	}

	// Check context
	if it.Ctx != nil {
		select {
		case <-it.Ctx.Done():
			it.err = it.Ctx.Err()
			return false
		default:
		}
	}

	// Check if we need to load more data
	if it.index >= len(it.current) {
		if !it.HasMoreData {
			return false // No more data
		}

		// Load next page
		it.Opts.Page = it.Page
		it.err = it.loadPage()
		if it.err != nil {
			return false
		}

		if len(it.current) == 0 {
			it.HasMoreData = false
			return false // No more data
		}

		it.index = 0
		it.Page++
	}

	// Advance index for next call
	it.index++
	return true
}

// Current returns the current item (call after Next() returns true)
func (it *SettingsIterator) Current() *commands.Setting {
	if it.index == 0 || it.index > len(it.current) {
		return nil
	}
	return it.current[it.index-1]
}

// HasMore returns true if there are more items to iterate
func (it *SettingsIterator) HasMore() bool {
	return it.HasMoreData || it.index < len(it.current)
}

// Err returns the last error encountered
func (it *SettingsIterator) Err() error {
	return it.err
}

// Total returns the total number of items (if available)
func (it *SettingsIterator) Total() int {
	return it.total
}

// Close cleans up the iterator
func (it *SettingsIterator) Close() {
	it.HasMoreData = false
	it.current = nil
}
