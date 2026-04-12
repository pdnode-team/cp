<!-- 例如 Login.svelte -->
<script lang="ts">
	import pb from '$lib/pocketbase';

	let errorText = $state('');

	let email = $state('');
	let password = $state('');

	const handleLogin = async (e: Event) => {
		errorText = '';
		e.preventDefault();
		if (!email || !password) {
			errorText = 'Please fill in all fields';
			return;
		}

		try {
			await pb.collection('users').authWithPassword(email, password);

			window.location.pathname = '/';
		} catch (err: any) {
			const errorData = err.data?.data || {};

			if (Object.keys(errorData).length === 0 && err.status == 400){
				errorText = err.data.message
				return
			}

			// 获取第一个字段的错误消息
			// 例如："passwordConfirm: Values don't match."
			const firstKey = Object.keys(errorData)[0];
			if (firstKey) {
				const friendlyMessage = `${firstKey}: ${errorData[firstKey].message}`;
				console.error(friendlyMessage);
				errorText = friendlyMessage;
				return;
			}
		}
	};
</script>

<div class="flex min-h-[70vh] flex-col items-center justify-center px-4 py-12">
	<div class="card w-full max-w-md border border-base-200 bg-base-100 shadow-2xl">
		<div class="card-body p-8">
			<h2 class="mb-2 text-center text-3xl font-bold">Login</h2>
			<p class="mb-6 text-center text-base-content/60">
				Enter your email and password below to log in
			</p>
			<div class="flex flex-col gap-4">
				<label class="form-control w-full">
					<div class="label"><span class="label-text font-medium">Email</span></div>
					<input
						type="email"
						name="email"
						placeholder="mail@example.com"
						class="input-bordered input w-full"
						autocomplete="email"
						required
						bind:value={email}
					/>
				</label>

				<label class="form-control w-full">
					<div class="label"><span class="label-text font-medium">Password</span></div>
					<input
						type="password"
						name="password"
						placeholder="••••••••"
						class="input-bordered input w-full"
						required
						bind:value={password}
					/>
				</label>

				<p class="text-red-500">{errorText}</p>

				<div class="form-control mt-6">
					<button onclick={handleLogin} type="submit" class="btn w-full text-lg btn-primary"
						>Log in</button
					>
				</div>

				<div class="mt-4 text-center text-sm">
					Don't have an account? <a href="/register" class="link font-medium link-primary"
						>Register</a
					>
				</div>
			</div>
		</div>
	</div>
</div>
