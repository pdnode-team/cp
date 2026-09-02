<script lang="ts">
	import './layout.css'

	let { children } = $props()
	let user: any = $state()

	import pb from '$lib/pocketbase'
	import { goto } from '$app/navigation'
	import { onMount } from 'svelte'

	function reloadLoginStatus() {
		if (pb.authStore.isValid) {
			user = pb.authStore.record
		} else {
			user = null
		}
		
	}
	pb.authStore.onChange((_, record) => {
		reloadLoginStatus()
	}, true)

	onMount(() => {
		reloadLoginStatus()
	})
</script>

<!-- 顶部导航栏 -->
<header
	class="navbar mx-auto max-w-7xl border-b border-base-200 bg-base-100 px-4 shadow-sm md:px-8"
>
	<div class="flex-1">
		<a href="/" class="flex items-center gap-2 text-lg font-bold tracking-tight transition-colors hover:text-primary sm:text-xl">
			<span>Pdnode CP<span class="hidden sm:inline"> Website</span></span>
			<div class="badge badge-soft text-xs">v0.5.1</div>
		</a>
	</div>

	<!-- 桌面端导航（md 及以上） -->
	<div class="hidden flex-none items-center gap-4 md:flex">
		<a href="/about" class="btn text-base-content/80 btn-ghost btn-sm">About</a>

		{#if user}
			<a href="/explore" class="btn text-base-content/80 btn-ghost btn-sm">Explore</a>
			<div class="dropdown dropdown-end">
				<button class="placeholder btn avatar btn-circle btn-ghost">
					<div
						class="flex w-10 items-center justify-center rounded-full bg-primary text-primary-content"
					>
						<span class="text-lg font-bold leading-none">
							{user.name?.trim().at(0) ?? '?'}
						</span>
					</div>
				</button>
				<ul class="dropdown-content menu z-50 mt-3 flex w-52 flex-col gap-2 rounded-box bg-base-300 p-2 shadow menu-sm">
					<li>
						<button
							onclick={() => {
								pb.authStore.clear()
								window.location.reload()
							}}
							class="w-full text-left text-error">Logout</button
						>
					</li>
					<li>
						<button onclick={() => goto("/mystuff")}>My Stuff</button>
					</li>
				</ul>
			</div>
		{:else}
			<a href="/register" class="btn btn-outline btn-sm">Register</a>
			<a href="/login" class="btn btn-primary btn-sm">Login</a>
		{/if}
	</div>

	<!-- 移动端导航（< md） -->
	<div class="flex flex-none items-center gap-2 md:hidden">
		{#if user}
			<div class="placeholder avatar">
				<div class="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-xs font-bold text-primary-content">
					{user.name?.trim().at(0) ?? '?'}
				</div>
			</div>
		{/if}

		<div class="dropdown dropdown-end">
			<button tabindex="0" class="btn btn-circle btn-ghost btn-sm" aria-label="Open menu">
				<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
				</svg>
			</button>
			<ul class="dropdown-content menu z-50 mt-3 w-56 rounded-box bg-base-200 p-3 shadow-xl">
				{#if user}
					<li class="menu-title px-2 py-1 text-xs opacity-60">
						Account: {user.name || user.email}
					</li>
					<li><a href="/explore">Explore</a></li>
					<li><a href="/create">Create CP</a></li>
					<li><a href="/mystuff">My Stuff</a></li>
					<li><a href="/about">About</a></li>
					<div class="divider my-1"></div>
					<li>
						<button
							onclick={() => {
								pb.authStore.clear()
								window.location.reload()
							}}
							class="text-error"
						>Logout</button>
					</li>
				{:else}
					<li><a href="/about">About</a></li>
					<div class="divider my-1"></div>
					<li class="mb-1"><a href="/login" class="btn btn-primary btn-sm text-primary-content">Login</a></li>
					<li><a href="/register" class="btn btn-outline btn-sm">Register</a></li>
				{/if}
			</ul>
		</div>
	</div>
</header>

<!-- 主体内容插槽 -->
<main class="mx-auto min-h-[calc(100vh-65px)] max-w-7xl border-base-200 bg-base-100 md:border-x">
	{@render children()}
</main>
