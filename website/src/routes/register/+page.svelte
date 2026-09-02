<script lang="ts">
	import { goto } from '$app/navigation';
	import pb from '$lib/pocketbase';

	let email = $state('');
	let password = $state('');
	let passwordConfirm = $state('');
	let name = $state('');
	let errorText = $state('');

	const handleRegister = async (e: Event) => {
		errorText = '';
		e.preventDefault();
		if (!email || !password || !name || !passwordConfirm) {
			errorText = 'Please fill in all fields';
			return;
		}else if (password !== passwordConfirm) {
			errorText = "The confirmation password does not match the password."

			return
		}

		try {
			await pb.collection('users').create({
				email,
				password,
				name,
				passwordConfirm
			});

			goto('/login')

		} catch (err: any) {
			const errorData = err.data?.data || {};

			// 获取第一个字段的错误消息
			// 例如："passwordConfirm: Values don't match."
			const firstKey = Object.keys(errorData)[0];
			if (firstKey) {
				const friendlyMessage = `${firstKey}: ${errorData[firstKey].message}`;
				console.error(friendlyMessage);
				errorText = friendlyMessage
				return
			}
		}
	};
</script>

<div class="flex min-h-[70vh] flex-col items-center justify-center px-4 py-8 sm:py-12">
	<div class="card w-full max-w-md border border-base-200 bg-base-100 shadow-2xl">
		<div class="card-body p-6 sm:p-8">
			<h2 class="mb-2 text-center text-3xl font-bold">Register</h2>
			<p class="mb-6 text-center text-base-content/60">
				Enter your email and password below to register
			</p>
			<div class="flex flex-col gap-4">
				<label class="form-control w-full">
					<div class="label"><span class="label-text font-medium">Name</span></div>
					<input
						type="text"
						name="name"
						placeholder="John Joe"
						class="input-bordered input w-full"
						autocomplete="name"
						required
						bind:value={name}
					/>
				</label>

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

				<label class="form-control w-full">
					<div class="label"><span class="label-text font-medium">Confirm Password</span></div>
					<input
						type="password"
						name="password"
						placeholder="••••••••"
						class="input-bordered input w-full"
						required
						bind:value={passwordConfirm}
					/>
				</label>

				<p class="text-red-500">{errorText}</p>

				<!-- 因为父级有了 gap-4 (1rem)，这里的 mt-6 (1.5rem) 会叠加，按钮会和上面隔得更远一点，显得有呼吸感。如果你不需要这么远，把 mt-6 删掉即可 -->
				<div class="form-control mt-6">
					<button onclick={handleRegister} type="submit" class="btn w-full text-lg btn-primary"
						>Register</button
					>
				</div>

				<div class="mt-4 text-center text-sm">
					Have an account? <a href="/login" class="link font-medium link-primary">Log in</a>
				</div>
			</div>
		</div>
	</div>
</div>
